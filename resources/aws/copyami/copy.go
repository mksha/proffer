package copyami

import (
	"fmt"
	"sync"

	awscommon "example.com/proffer/resources/aws/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type SrcAmiInfo struct {
	Region  *string
	Filters []*ec2.Filter
	Image   *ec2.Image
}

type TargetInfo struct {
	Regions  []*string
	CopyTags bool
	Tags []*ec2.Tag
}

var wg sync.WaitGroup

// type AmiDoesNotExist struct {
// 	Filters    []*ec2.Filter
// 	Region     *string
// 	StatusCode string
// }

// func (a AmiDoesNotExist) Error() string {
// 	a.StatusCode = "AmiDoesNotExist"
// 	return fmt.Sprintf("%s: AMI does not exist in region %s with filters %v ", a.StatusCode, *a.Region, a.Filters)
// }

func isError(err error) (bool, error) {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "RequestExpired":
				return true, fmt.Errorf("%v: Provided credential has expired", aerr.Code())
			default:
				return true, fmt.Errorf("%s", aerr.Error())
			}
		}
		return true, err
	}
	return false, nil
}

func isAmiExist(sess *session.Session, filters []*ec2.Filter) (bool, error) {
	svc := ec2.New(sess)
	input := &ec2.DescribeImagesInput{
		Filters: filters,
	}

	result, err := svc.DescribeImages(input)
	if ok, err := isError(err); ok {
		return false, err
	}

	images := result.Images

	if len(images) == 0 {
		return false, nil
	}
	return true, nil
}

func getAmiInfo(sess *session.Session, filters []*ec2.Filter) ([]*ec2.Image, error) {
	if ok, err := isAmiExist(sess, filters); !ok {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("UnableToGetAmiInfo: AMI doesnot exist in Region %s with Filters %s ", *sess.Config.Region, filters)
	}

	svc := ec2.New(sess)
	input := &ec2.DescribeImagesInput{
		Filters: filters,
	}

	result, err := svc.DescribeImages(input)
	if ok, err := isError(err); ok {
		return nil, err
	}

	images := result.Images
	return images, nil
}

func createEc2Tags(sess *session.Session, resources []*string, tags []*ec2.Tag) error {
	svc := ec2.New(sess)
	input := &ec2.CreateTagsInput{
		Resources: resources,
		Tags: tags,
	}

	_, err := svc.CreateTags(input)
	if err != nil {
		return err
	}

	return nil
}


func formEc2Tags(tags map[*string]*string) []*ec2.Tag {
	ec2Tags := make([]*ec2.Tag, len(tags))

	for key, value := range tags {
		ec2Tags = append(ec2Tags, &ec2.Tag{Key:key,Value:value})
	}

	return ec2Tags
}


func copyImage(sess *session.Session, sai SrcAmiInfo, tags []*ec2.Tag, errMap map[string]error) {
	defer wg.Done()
	filters := []*ec2.Filter{
		{
			Name: aws.String("name"),
			Values: []*string{
				sai.Image.Name,
			},
		},
	}
	ok, err := isAmiExist(sess, filters)
	if ok {
		clogger.Warnf("AMI %s Already Exist In Region %s", *sai.Image.Name, *sess.Config.Region)
		return
	} else {
		if err != nil {
			errMap[*sess.Config.Region] = err
			return
		}
	}

	clogger.Infof("Started Copying AMI In Region: %s ...", *sess.Config.Region)

	svc := ec2.New(sess)
	input := &ec2.CopyImageInput{
		Description:   sai.Image.Description,
		Name:          sai.Image.Name,
		SourceImageId: sai.Image.ImageId,
		SourceRegion:  sai.Region,
	}

	result, err := svc.CopyImage(input)

	if ok, err := isError(err); ok {
		errMap[*sess.Config.Region] = err
		return
	}

	err = svc.WaitUntilImageAvailable(&ec2.DescribeImagesInput{ImageIds: []*string{result.ImageId}})
	if err != nil {
		errMap[*sess.Config.Region] = err
		return
	}

	clogger.Infof("Copied AMI In Region: %s , New AMI Id Is: %s", *sess.Config.Region, *result.ImageId)

	if len(tags) == 0 {
		clogger.Debug("No Tags To Add Or Create")
		return
	}

	clogger.Debugf("Adding Following Tags to AMI %s",*result.ImageId)
	clogger.Debug(tags)

	if err := createEc2Tags(sess, []*string{result.ImageId}, tags); err != nil {
		errMap[*sess.Config.Region] = err
		return
	}

	clogger.Infof("Tags Have Copied/Added To AMI : %s , In Region: %s", *result.ImageId, *sess.Config.Region)
}

func copyAmi(srcAmiInfo SrcAmiInfo, targetInfo TargetInfo) error {

	sess, err := awscommon.GetAwsSessWithDefaultCreds()
	if err != nil {
		return err
	}

	sess.Config.Region = srcAmiInfo.Region
	images, err := getAmiInfo(sess, srcAmiInfo.Filters)

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
		go copyImage(sess.Copy(&aws.Config{Region: targetRegion}), srcAmiInfo, targetInfo.Tags, errMap)
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
