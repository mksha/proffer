package common

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func IsError(err error) (bool, error) {
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

func IsAmiExist(sess *session.Session, filters []*ec2.Filter) (bool, error) {
	svc := ec2.New(sess)
	input := &ec2.DescribeImagesInput{
		Filters: filters,
	}

	result, err := svc.DescribeImages(input)
	if ok, err := IsError(err); ok {
		return false, err
	}

	images := result.Images

	if len(images) == 0 {
		return false, nil
	}
	return true, nil
}

func GetAmiInfo(sess *session.Session, filters []*ec2.Filter) ([]*ec2.Image, error) {
	if ok, err := IsAmiExist(sess, filters); !ok {
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("UnableToGetAmiInfo: AMI doesnot exist in Region %s with Filters %v ", *sess.Config.Region, filters)
	}

	svc := ec2.New(sess)
	input := &ec2.DescribeImagesInput{
		Filters: filters,
	}

	result, err := svc.DescribeImages(input)
	if ok, err := IsError(err); ok {
		return nil, err
	}

	images := result.Images
	return images, nil
}

func CreateEc2Tags(sess *session.Session, resources []*string, tags []*ec2.Tag) error {
	svc := ec2.New(sess)
	input := &ec2.CreateTagsInput{
		Resources: resources,
		Tags:      tags,
	}

	_, err := svc.CreateTags(input)
	if err != nil {
		return err
	}

	return nil
}

func FormEc2Tags(tags map[*string]*string) []*ec2.Tag {
	ec2Tags := make([]*ec2.Tag, 0)

	for key, value := range tags {
		ec2Tags = append(ec2Tags, &ec2.Tag{Key: key, Value: value})
	}

	return ec2Tags
}
