package shareami

import (
	// "fmt"
	"fmt"
	"sync"

	awscommon "example.com/proffer/resources/aws/common"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var wg sync.WaitGroup
var innerWg sync.WaitGroup

type (
	RegionErrMap        map[string]error
	AccountRegionErrMap map[string]RegionErrMap
)

func addPermissionsforTargetAccounts(sai SrcAmiInfo, region *string, regionErrMap RegionErrMap, accountFlagMap AccountFlagMap) {
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

	for account, flag := range accountFlagMap {
		rawAccounts = append(rawAccounts, *account)

		if flag.AddCVP {
			flaggedRawAccounts = append(flaggedRawAccounts, *account)
		}
	}

	clogger.Infof("\t  With Account(s): %v", rawAccounts)
	clogger.Infof("\t  In Region: %s", *region)
	clogger.Info("")

	launchPermissions := make([]*ec2.LaunchPermission, 0)
	createVolumePermissions := make([]*ec2.CreateVolumePermission, 0)

	for targetAccount, flag := range accountFlagMap {
		launchPermission := &ec2.LaunchPermission{UserId: targetAccount}
		launchPermissions = append(launchPermissions, launchPermission)

		if flag.AddCVP {
			clogger.Debugf("Found 'addCreateVolumeFlag' Flag Set to 'true' For Target Account: %s", *targetAccount)
			clogger.Debug("\t  Will Add 'CreateVolumePermission' To")
			clogger.Debugf("\t  Source AMI: %s In Region: %s For Target Account: %s", *image.Name, *region, *targetAccount)
			clogger.Debug("")

			createVolumePermission := &ec2.CreateVolumePermission{UserId: targetAccount}
			createVolumePermissions = append(createVolumePermissions, createVolumePermission)
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
	flag := Flag{AddCVP: m.AddCVP}
	accountFlagMap := AccountFlagMap{m.AccountID: flag}

	for _, targetRegion := range m.Regions {
		wg.Add(1)

		go addPermissionsforTargetAccounts(sai, targetRegion, regionErrMap, accountFlagMap)
	}

	wg.Wait()

	if len(regionErrMap) != 0 {
		accountRegionErrMap[*m.AccountID] = regionErrMap
	}
}

func (t Target) shareAmiWithCommonRegions(sai SrcAmiInfo, commonAccountRegionErrMap AccountRegionErrMap) {
	regionErrMap := RegionErrMap{}
	accountFlagMap := t.getTargetAccountFlagMap()

	for _, targetRegion := range t.CommonRegions {
		wg.Add(1)

		go addPermissionsforTargetAccounts(sai, targetRegion, regionErrMap, accountFlagMap)
	}

	wg.Wait()

	if len(regionErrMap) != 0 {
		for targetAccount, _ := range accountFlagMap {
			commonAccountRegionErrMap[*targetAccount] = regionErrMap
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

	return nil
}
