package shareami

import (
	"fmt"

	awscommon "example.com/proffer/resources/aws/common"
	// "github.com/aws/aws-sdk-go/service/ec2"
)


func (c *Config) shareAmi() {

}

func (c *Config) apply() error {
	sess, err := awscommon.GetAwsSession(c.SrcAmiInfo.CredsInfo)
	if err != nil {
		return err
	}

	accountInfo, err := awscommon.GetAccountInfo(sess)
	if err != nil {
		return err
	}

	c.SrcAmiInfo.AccountID = accountInfo.Account
	sess.Config.Region = c.SrcAmiInfo.Region
	images, err := awscommon.GetAmiInfo(sess, c.SrcAmiInfo.Filters)

	if err != nil {
		return err
	}

	c.SrcAmiInfo.Image = images[0]

	for _, AccountRegionMapping := range c.Target.ModAccountRegionMappingList {
		fmt.Println(AccountRegionMapping)
	}

	return nil
}
