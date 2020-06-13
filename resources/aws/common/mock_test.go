package common

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

// Define a mock struct for ec2 to use in unit tests.
type mockEC2Client struct {
	ec2iface.EC2API
}

// Define a mock struct for sts to use in unit tests.
type mockSTSClient struct {
	stsiface.STSAPI
	Error error
}

// Define a mock struct for iam to use in unit tests.
type mockIAMClient struct {
	iamiface.IAMAPI
	Error error
}

func (m *mockEC2Client) DescribeImages(input *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	// Check that required inputs exist
	if input.Filters == nil || len(input.Filters) == 0 {
		return nil, errors.New("DescribeImagesInput.Filters is nil or an empty slice")
	}

	var resp *ec2.DescribeImagesOutput

	for _, filter := range input.Filters {
		if *filter.Name == "unknown" {
			resp = &ec2.DescribeImagesOutput{
				Images: []*ec2.Image{},
			}

			return resp, nil
		}
	}

	resp = &ec2.DescribeImagesOutput{
		Images: []*ec2.Image{{
			ImageId: aws.String("test-image-id"),
		}},
	}

	return resp, nil
}

func (m *mockEC2Client) CreateTags(input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	// Check that required inputs exist
	if input.Resources == nil || len(input.Resources) == 0 {
		return nil, errors.New("CreateTagsInput.Resources is nil or an empty slice")
	}

	if input.Tags == nil || len(input.Tags) == 0 {
		return nil, errors.New("CreateTagsInput.Tags is nil or an empty slice")
	}

	if input.Tags[0].Key == nil || *input.Tags[0].Key == "" || input.Tags[0].Value == nil || *input.Tags[0].Value == "" {
		return nil, errors.New("CreateTagsInput.Tags[0].Tag or CreateTagsInput.Tags[0].Value is nil or an empty string")
	}

	resp := &ec2.CreateTagsOutput{}

	return resp, nil
}

func (m *mockSTSClient) GetCallerIdentity(input *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	resp := &sts.GetCallerIdentityOutput{
		Account: aws.String("123456789012"),
		Arn:     aws.String("arn:aws::iam:123456789012:role/test"),
		UserId:  aws.String("123@testuser"),
	}

	return resp, m.Error
}

func (m *mockIAMClient) ListAccountAliases(input *iam.ListAccountAliasesInput) (*iam.ListAccountAliasesOutput, error) {
	resp := &iam.ListAccountAliasesOutput{
		AccountAliases: []*string{
			aws.String("test-account"),
		},
		IsTruncated: aws.Bool(false),
	}

	return resp, m.Error
}
