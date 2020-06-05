package copyami

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	awscommon "github.com/proffer/resources/aws/common"
	"github.com/aws/aws-sdk-go/service/sts"
)

type TargetInfo struct {
	Regions  []*string
	CopyTags bool
	Tags     []*ec2.Tag
}

var wg sync.WaitGroup

// It copies the given source ami to target regions.
func copyAmi(sess *session.Session, sai SrcAmiInfo, tags []*ec2.Tag, errMap map[string]error) {
	defer wg.Done()

	svc := ec2.New(sess)

	ok, err := awscommon.IsAmiExist(svc, sai.Filters)
	if ok {
		clogger.Warnf("AMI %s Already Exist In Account %s In Region %s", *sai.Image.Name, *sai.AccountID, *sess.Config.Region)
		return
	} else if err != nil {
		errMap[*sess.Config.Region] = err
		return
	}

	clogger.Infof("Started Copying AMI In Account: %s Region: %s ...", *sai.AccountID, *sess.Config.Region)

	input := &ec2.CopyImageInput{
		Description:   sai.Image.Description,
		Name:          sai.Image.Name,
		SourceImageId: sai.Image.ImageId,
		SourceRegion:  sai.Region,
	}

	result, err := svc.CopyImage(input)

	if ok, err := awscommon.IsError(err); ok {
		errMap[*sess.Config.Region] = err
		return
	}

	err = svc.WaitUntilImageAvailable(&ec2.DescribeImagesInput{ImageIds: []*string{result.ImageId}})
	if err != nil {
		errMap[*sess.Config.Region] = err
		return
	}

	clogger.Infof("Copied AMI In Account: %s In Region: %s , New AMI Id Is: %s", *sai.AccountID, *sess.Config.Region, *result.ImageId)

	if len(tags) == 0 {
		clogger.Debug("No Tags To Add Or Create")
		return
	}

	clogger.Debugf("Adding Following Tags to AMI: %s", *result.ImageId)
	clogger.Debug(tags)

	if err := awscommon.CreateEc2Tags(svc, []*string{result.ImageId}, tags); err != nil {
		errMap[*sess.Config.Region] = err
		return
	}

	clogger.Infof("Tags Have Copied/Added To AMI : %s , In Region: %s", *result.ImageId, *sess.Config.Region)
}

// It applies the configuration for resources of kind aws-copyami.
func apply(srcAmiInfo SrcAmiInfo, targetInfo TargetInfo) error {
	sess, err := awscommon.GetAwsSession(srcAmiInfo.CredsInfo)
	if err != nil {
		return err
	}

	svc := sts.New(sess)
	accountInfo, err := awscommon.GetAccountInfo(svc)
	if err != nil {
		return err
	}

	srcAmiInfo.AccountID = accountInfo.Account
	sess.Config.Region = srcAmiInfo.Region
	ci := awscommon.AwsClientInfo{
		SVC:    ec2.New(sess),
		Region: sess.Config.Region,
	}
	images, err := awscommon.GetAmiInfo(ci, srcAmiInfo.Filters)

	if err != nil {
		return err
	}

	srcAmiInfo.Image = images[0]
	errMap := map[string]error{}

	if targetInfo.CopyTags {
		targetInfo.Tags = append(targetInfo.Tags, srcAmiInfo.Image.Tags...)
	}

	for _, targetRegion := range targetInfo.Regions {
		wg.Add(1)

		go copyAmi(sess.Copy(&aws.Config{Region: targetRegion}), srcAmiInfo, targetInfo.Tags, errMap)
	}

	wg.Wait()

	if len(errMap) != 0 {
		clogger.Error("AMI Copy Operation Failed For following Regions:")

		for region, err := range errMap {
			clogger.Errorf("%s:\n", region)
			clogger.Errorf("\tReason: [%s] ", err)
		}

		return fmt.Errorf("Failed")
	}

	return nil
}
