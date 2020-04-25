package copyami

import (
	"fmt"
	"log"
	"runtime"
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
	Accounts []*string
}

var wg sync.WaitGroup

// type AmiDoesNotExist struct {
// 	Filters []*ec2.Filter
// 	Region  *string
// 	StatusCode string
// }

// func (a AmiDoesNotExist) Error() string {
// 	return fmt.Sprintf("AmiDoesNotExist: AMI does not exist in region %s with filters %v ", *a.Region, a.Filters)
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
		return false, fmt.Errorf(" AMIDoesNotExist: No AMI Matched With Given Filters")
	}
	return true, nil
}

func getAmiInfo(sess *session.Session, filters []*ec2.Filter) ([]*ec2.Image, error) {
	if ok, err := isAmiExist(sess, filters); !ok {
		return nil, err
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

func copyImage(sess *session.Session, sai SrcAmiInfo) {
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
		log.Printf(" AMI %s Already Exist In Region %s", *sai.Image.Name, *sess.Config.Region)
		return
	} else {
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Printf(" Start Copying AMI In Region %s ...", *sess.Config.Region)

	svc := ec2.New(sess)
	input := &ec2.CopyImageInput{
		Description:   sai.Image.Description,
		Name:          sai.Image.Name,
		SourceImageId: sai.Image.ImageId,
		SourceRegion:  sai.Region,
	}

	result, err := svc.CopyImage(input)

	if ok, err := isError(err); ok {
		log.Fatalln(err)
	}

	err = svc.WaitUntilImageAvailable(&ec2.DescribeImagesInput{ImageIds: []*string{result.ImageId}})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf(" Copied AMI In Region: %s , New AMI Id Is : %s ", sess.Config.Region, &result.ImageId)
}

func copyAmi(srcAmiInfo SrcAmiInfo, targetInfo TargetInfo) {
	log.Println(" ******************** Start: Copy AMI Operation ************************************")

	sess := awscommon.GetAwsSessWithDefaultCreds()
	sess.Config.Region = srcAmiInfo.Region
	images, err := getAmiInfo(sess, srcAmiInfo.Filters)

	if err != nil {
		log.Fatalln(err)
	}

	srcAmiInfo.Image = images[0]
	wg.Add(len(targetInfo.Regions))
	for _, targetRegion := range targetInfo.Regions {
		go copyImage(sess.Copy(&aws.Config{Region: targetRegion}), srcAmiInfo)
	}

	fmt.Println("GoRutines:", runtime.NumGoroutine())
	wg.Wait()
	fmt.Println("GoRutines:", runtime.NumGoroutine())
	defer log.Println(" ******************** END: Copy AMI Operation **************************************")
}
