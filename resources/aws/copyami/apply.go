package copyami

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	awscommon "github.com/proffer/resources/aws/common"
)

// TargetInfo represents the target data in go struct form.
type TargetInfo struct {
	Regions  []*string
	CopyTags bool
	Tags     []*ec2.Tag
}

var wg sync.WaitGroup

// It copies the given source ami to target regions.
func (r *Resource) copyAmi(sess *session.Session, sai SrcAmiInfo, tags []*ec2.Tag, errMap map[string]error) {
	defer wg.Done()

	svc := ec2.New(sess)
	amiMeta := AmiMeta{Name: sai.Image.Name}

	defer func() {
		r.Record.TargetImages[sess.Config.Region] = amiMeta
	}()

	targetAmiFilters := []*ec2.Filter{{
		Name:   aws.String("name"),
		Values: []*string{amiMeta.Name}}}

	ok, err := awscommon.IsAmiExist(svc, targetAmiFilters)
	if ok {
		ci := awscommon.AwsClientInfo{
			SVC:    ec2.New(sess),
			Region: sess.Config.Region,
		}
		images, _ := awscommon.GetAmiInfo(ci, targetAmiFilters)

		clogger.Warnf("AMI %s(%s) Already Exist In Account %s In Region %s", *sai.Image.Name, *images[0].ImageId, *sai.AccountID, *sess.Config.Region)
		amiMeta.ID = images[0].ImageId

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

	amiMeta.ID = result.ImageId

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
func (r *Resource) apply(srcAmiInfo SrcAmiInfo, targetInfo TargetInfo) error {
	sess, err := awscommon.GetAwsSession(srcAmiInfo.CredsInfo)
	if err != nil {
		return err
	}

	svc := sts.New(sess)

	callerInfo, err := awscommon.GetCallerInfo(svc)
	if err != nil {
		return err
	}

	iamSVC := iam.New(sess)

	accountAlias, err := awscommon.GetAccountAlias(iamSVC)
	if err != nil {
		return err
	}

	r.Record.AccountMeta = AccountMeta{
		ID:    callerInfo.Account,
		Alias: accountAlias,
	}

	srcAmiInfo.AccountID = callerInfo.Account
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
	r.Record.SrcImage.Region = srcAmiInfo.Region
	r.Record.SrcImage.Name = images[0].Name
	r.Record.SrcImage.ID = images[0].ImageId
	r.Record.TargetImages = make(map[*string]AmiMeta)
	errMap := map[string]error{}

	if targetInfo.CopyTags {
		targetInfo.Tags = append(targetInfo.Tags, srcAmiInfo.Image.Tags...)
	}

	for _, targetRegion := range targetInfo.Regions {
		wg.Add(1)

		go r.copyAmi(sess.Copy(&aws.Config{Region: targetRegion}), srcAmiInfo, targetInfo.Tags, errMap)
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
