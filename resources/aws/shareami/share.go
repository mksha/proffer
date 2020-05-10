package shareami

import (
	// "fmt"
	"fmt"
	"sync"

	awscommon "example.com/proffer/resources/aws/common"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	wg         sync.WaitGroup
	innerWg    sync.WaitGroup
	copyTagsWg sync.WaitGroup
	triggerWg  sync.WaitGroup
)

type (
	RegionErrMap        map[string]error
	AccountRegionErrMap map[string]RegionErrMap
)

func (c Config) triggerCopyTagsAcrossAccounts() error {
	accountRegionErrMap := AccountRegionErrMap{}

	for _, accountRegionMapping := range c.Target.ModAccountRegionMappingList {
		if accountRegionMapping.CopyTags {
			triggerWg.Add(1)

			go copyTagsAcrossAccount(c.SrcAmiInfo, accountRegionMapping, c.Target.CommonRegions, accountRegionErrMap)
		}
	}

	triggerWg.Wait()

	if len(accountRegionErrMap) != 0 {
		clogger.Error("'Copy Tags Across Accounts' Operation Failed For Following Accounts:")

		for account, regionErrMap := range accountRegionErrMap {
			clogger.Errorf("\t- Account: %s", account)

			if v, ok := regionErrMap["*"]; ok {
				clogger.Error("\t  - Regions: All")
				clogger.Errorf("\t    Reason: [%s] ", v)
			} else {
				clogger.Error("\t  - Regions:")
				for region, err := range regionErrMap {
					clogger.Errorf("\t    - Region: %s", region)
					clogger.Errorf("\t      Reason: [%s] ", err)
				}
			}
		}

		return fmt.Errorf("exiting.")
	}

	return nil
}

func copyTagsAcrossAccount(sai SrcAmiInfo, accountRegionMapping AccountRegionMapping, commonRegions []*string, accountRegionErrMap AccountRegionErrMap) {
	defer triggerWg.Done()

	regionErrMap := RegionErrMap{}

	sess, err := awscommon.GetAwsSession(accountRegionMapping.CredsInfo)
	if err != nil {
		accountRegionErrMap[*accountRegionMapping.AccountID] = RegionErrMap{"*": err}
		return
	}

	targetRegions := make([]*string, 0)
	targetRegions = append(targetRegions, accountRegionMapping.Regions...)
	targetRegions = append(targetRegions, commonRegions...)
	rawTargetRegions := make([]string, 0)

	for _, targetRegion := range targetRegions {
		rawTargetRegions = append(rawTargetRegions, *targetRegion)
		copyTagsWg.Add(1)

		go addTagsToTargetAmi(sess.Copy(&aws.Config{Region: targetRegion}), sai, accountRegionMapping.Tags, regionErrMap)
	}

	copyTagsWg.Wait()

	if len(regionErrMap) != 0 {
		accountRegionErrMap[*accountRegionMapping.AccountID] = regionErrMap
	}

	clogger.Success("Successfully Added/Copied Tags To AMI(s)")
	clogger.Infof("\t  In Target Account: %s", *accountRegionMapping.AccountID)
	clogger.Infof("\t  In Regions: %v", rawTargetRegions)
	clogger.Info("")

}

func addTagsToTargetAmi(sess *session.Session, sai SrcAmiInfo, tags []*ec2.Tag, regionErrMap RegionErrMap) {
	defer copyTagsWg.Done()

	region := sess.Config.Region
	amiID := sai.RegionAmiErrMap[region].Ami.ImageId
	snapshotID := sai.RegionAmiErrMap[region].Ami.BlockDeviceMappings[0].Ebs.SnapshotId
	tags = append(tags, sai.RegionAmiErrMap[region].Ami.Tags...)

	if err := awscommon.CreateEc2Tags(sess, []*string{amiID, snapshotID}, tags); err != nil {
		regionErrMap[*region] = err
		return
	}
}

func addPermissionsforTargetAccounts(sai SrcAmiInfo, region *string, regionErrMap RegionErrMap, accountRegionMappingList []AccountRegionMapping) {
	defer wg.Done()

	sess, err := awscommon.GetAwsSession(sai.CredsInfo)
	if err != nil {
		regionErrMap[*region] = err
		return
	}

	sess.Config.Region = region

	amiInfo, ok := sai.RegionAmiErrMap[region]
	if !ok {
		regionErrMap[*region] = amiInfo.Error
		return
	}

	image := amiInfo.Ami
	svc := ec2.New(sess)

	// modify source image launch permission to add target account
	clogger.Successf("Started Sharing AMI: %s", *image.Name)

	rawAccounts := make([]string, 0)
	flaggedRawAccounts := make([]string, 0)

	for _, accountRegionMapping := range accountRegionMappingList {
		rawAccounts = append(rawAccounts, *accountRegionMapping.AccountID)

		if accountRegionMapping.AddCVP {
			flaggedRawAccounts = append(flaggedRawAccounts, *accountRegionMapping.AccountID)
		}
	}

	clogger.Infof("\t  With Account(s): %v", rawAccounts)
	clogger.Infof("\t  In Region: %s", *region)
	clogger.Info("")

	launchPermissions := make([]*ec2.LaunchPermission, 0)
	createVolumePermissions := make([]*ec2.CreateVolumePermission, 0)

	for _, accountRegionMapping := range accountRegionMappingList {
		launchPermission := &ec2.LaunchPermission{UserId: accountRegionMapping.AccountID}
		launchPermissions = append(launchPermissions, launchPermission)

		if accountRegionMapping.AddCVP {
			clogger.Debugf("Found 'addCreateVolumeFlag' Flag Set to 'true' For Target Account: %s", *accountRegionMapping.AccountID)
			clogger.Debug("\t  Will Add 'CreateVolumePermission' To")
			clogger.Debugf("\t  Source AMI: %s In Region: %s For Target Account: %s", *image.Name, *region, *accountRegionMapping.AccountID)
			clogger.Debug("")

			createVolumePermission := &ec2.CreateVolumePermission{UserId: accountRegionMapping.AccountID}
			createVolumePermissions = append(createVolumePermissions, createVolumePermission)
		}

		if accountRegionMapping.CopyTags {
			accountRegionMapping.Tags = append(accountRegionMapping.Tags, image.Tags...)
		}
	}

	modifyImageAttributeInput := &ec2.ModifyImageAttributeInput{
		ImageId: image.ImageId,
		LaunchPermission: &ec2.LaunchPermissionModifications{
			Add: launchPermissions,
		},
	}

	_, err = svc.ModifyImageAttribute(modifyImageAttributeInput)
	if err != nil {
		regionErrMap[*region] = err
		return
	}

	clogger.Debug("Added 'LaunchPermission' To")
	clogger.Debugf("\t  Source AMI: %s", *image.Name)
	clogger.Debugf("\t  In Source Account: %s ", *sai.AccountID)
	clogger.Debugf("\t  In Region: %s ", *region)
	clogger.Debugf("\t  For Target Account(s): %v", rawAccounts)
	clogger.Debug("")

	if len(createVolumePermissions) != 0 {
		modifySnapshotAttributeInput := &ec2.ModifySnapshotAttributeInput{
			SnapshotId: image.BlockDeviceMappings[0].Ebs.SnapshotId,
			CreateVolumePermission: &ec2.CreateVolumePermissionModifications{
				Add: createVolumePermissions,
			},
		}

		_, err := svc.ModifySnapshotAttribute(modifySnapshotAttributeInput)
		if err != nil {
			regionErrMap[*region] = err
			return
		}

		clogger.Debug("Added 'CreateVolumePermission' To")
		clogger.Debugf("\t  Source AMI: %s", *image.Name)
		clogger.Debugf("\t  In Source Account: %s ", *sai.AccountID)
		clogger.Debugf("\t  In Region: %s ", *region)
		clogger.Debugf("\t  For Target Account(s): %v", flaggedRawAccounts)
		clogger.Debug("")
	}

	clogger.Successf("Successfully Shared AMI: %s", *image.Name)
	clogger.Infof("\t  With Account(s): %v", rawAccounts)
	clogger.Infof("\t  In Region: %s", *region)
	clogger.Info("")

}

func (m AccountRegionMapping) shareAmiWithIndividualRegions(sai SrcAmiInfo, accountRegionErrMap AccountRegionErrMap) {
	defer innerWg.Done()
	regionErrMap := RegionErrMap{}

	for _, targetRegion := range m.Regions {
		wg.Add(1)

		go addPermissionsforTargetAccounts(sai, targetRegion, regionErrMap, []AccountRegionMapping{m})
	}

	wg.Wait()

	if len(regionErrMap) != 0 {
		accountRegionErrMap[*m.AccountID] = regionErrMap
	}
}

func (t Target) shareAmiWithCommonRegions(sai SrcAmiInfo, commonAccountRegionErrMap AccountRegionErrMap) {
	regionErrMap := RegionErrMap{}

	for _, targetRegion := range t.CommonRegions {
		wg.Add(1)

		go addPermissionsforTargetAccounts(sai, targetRegion, regionErrMap, t.ModAccountRegionMappingList)
	}

	wg.Wait()

	if len(regionErrMap) != 0 {
		for _, accountRegionMapping := range t.ModAccountRegionMappingList {
			commonAccountRegionErrMap[*accountRegionMapping.AccountID] = regionErrMap
		}
	}
}

func (c *Config) apply() error {
	accountRegionErrMap := AccountRegionErrMap{}
	commonAccountRegionErrMap := AccountRegionErrMap{}

	c.Target.shareAmiWithCommonRegions(c.SrcAmiInfo, commonAccountRegionErrMap)

	for _, AccountRegionMapping := range c.Target.ModAccountRegionMappingList {
		innerWg.Add(1)

		go AccountRegionMapping.shareAmiWithIndividualRegions(c.SrcAmiInfo, accountRegionErrMap)
	}

	innerWg.Wait()

	if len(accountRegionErrMap) != 0 || len(commonAccountRegionErrMap) != 0 {
		clogger.Error("AMI Share Operation Failed For Following Accounts:")

		if len(accountRegionErrMap) != 0 {
			for account, regionErrMap := range accountRegionErrMap {
				clogger.Errorf("\t- Account: %s", account)
				clogger.Error("\t  - Regions:")

				for region, err := range regionErrMap {
					clogger.Errorf("\t    - Region: %s", region)
					clogger.Errorf("\t      Reason: [%s] ", err)
				}
			}
		}

		if len(commonAccountRegionErrMap) != 0 {
			for account, commonRegionErrMap := range commonAccountRegionErrMap {
				clogger.Errorf("\t- Account: %s", account)
				clogger.Error("\t  - Regions:")

				for commonRegion, err := range commonRegionErrMap {
					clogger.Errorf("\t    - Region: %s", commonRegion)
					clogger.Errorf("\t      Reason: [%s] ", err)
				}
			}
		}

		return fmt.Errorf("exiting.")
	}

	if err := c.triggerCopyTagsAcrossAccounts(); err != nil {
		return err
	}

	return nil
}
