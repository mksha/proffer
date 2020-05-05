package shareami

import (
	// "fmt"
	"sync"

	awscommon "example.com/proffer/resources/aws/common"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var wg sync.WaitGroup

type (
	RegionErrMap        map[string]error
	AccountRegionErrMap map[string]RegionErrMap
)

func (m *AccountRegionMapping) addPermissionsforTargetAccount(sai awscommon.SrcAmiInfo, region *string, regionErrMap RegionErrMap) {
	defer wg.Done()

	// modify source image launch permission to add target account
	sess, err := awscommon.GetAwsSession(sai.CredsInfo)
	if err != nil {
		regionErrMap[*region] = err
		return
	}

	clogger.Infof("Getting Source AMI Info From Account: %s Region: %s", *sai.AccountID, *region)

	sess.Config.Region = region
	images, err := awscommon.GetAmiInfo(sess, sai.Filters)

	if err != nil {
		regionErrMap[*region] = err
		return
	}

	image := images[0]

	clogger.Infof("Found Source AMI: %s In", *image.Name)
	clogger.Infof("\t  Account: %s", *sai.AccountID)
	clogger.Infof("\t  Region: %s", *region)

	svc := ec2.New(sess)
	input := &ec2.ModifyImageAttributeInput{
		ImageId: image.ImageId,
		LaunchPermission: &ec2.LaunchPermissionModifications{
			Add: []*ec2.LaunchPermission{
				{
					UserId: m.AccountID,
				},
			},
		},
	}

	clogger.Infof("Started Sharing AMI: %s", *image.Name)
	clogger.Infof("\t  With Account: %s", *m.AccountID)
	clogger.Infof("\t  In Region: %s", *region)

	_, err = svc.ModifyImageAttribute(input)
	if err != nil {
		regionErrMap[*region] = err
		return
	}

	clogger.Infof("Successfully Shared AMI: %s", *image.Name)
	clogger.Infof("\t  With Account: %s", *m.AccountID)
	clogger.Infof("\t  In Region: %s", *region)
}

func (m *AccountRegionMapping) shareAmi(sai awscommon.SrcAmiInfo, accountRegionErrMap AccountRegionErrMap) {
	regionErrMap := RegionErrMap{}

	for _, targetRegion := range m.Regions {
		wg.Add(1)

		go m.addPermissionsforTargetAccount(sai, targetRegion, regionErrMap)
	}

	wg.Wait()

	if len(regionErrMap) != 0 {
		accountRegionErrMap[*m.AccountID] = regionErrMap
	}
}

func (c *Config) apply() error {
	sess, err := awscommon.GetAwsSession(c.SrcAmiInfo.CredsInfo)
	if err != nil {
		return err
	}

	sess.Config.Region = c.SrcAmiInfo.Region
	accountInfo, err := awscommon.GetAccountInfo(sess)
	if err != nil {
		return err
	}

	c.SrcAmiInfo.AccountID = accountInfo.Account

	accountRegionErrMap := AccountRegionErrMap{}

	for _, AccountRegionMapping := range c.Target.ModAccountRegionMappingList {
		AccountRegionMapping.shareAmi(c.SrcAmiInfo, accountRegionErrMap)
	}

	if len(accountRegionErrMap) != 0 {
		clogger.Error("AMI Share Operation Failed For Following Accounts:")

		for account, regionErrMap := range accountRegionErrMap {
			clogger.Errorf("\t- Account: %s", account)
			clogger.Error("\t  - Regions:")

			for region, err := range regionErrMap {
				clogger.Errorf("\t    - Region: %s", region)
				clogger.Errorf("\t      Reason: [%s] ", err)
			}
		}
	}

	return nil
}
