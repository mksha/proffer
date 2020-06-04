package common

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// Define a mock struct to use in unit tests.
type mockEC2Client struct {
	ec2iface.EC2API
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

func TestIsError(t *testing.T) {
	for n := range isErrorTestCases {
		tt := isErrorTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := IsError(tt.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAmiExist(t *testing.T) {
	msvc := &mockEC2Client{}

	for n := range isAmiExistTestCases {
		tt := isAmiExistTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := IsAmiExist(msvc, tt.filters)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsAmiExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsAmiExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAmiInfo(t *testing.T) {
	for n := range getAmiInfoTestCases {
		tt := getAmiInfoTestCases[n]
		ci := AwsClientInfo{
			SVC:    &mockEC2Client{},
			Region: aws.String("us-east-1"),
		}

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetAmiInfo(ci, tt.filters)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAmiInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAmiInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateEc2Tags(t *testing.T) {
	for n := range createEc2TagsTestCases {
		msvc := &mockEC2Client{}

		tt := createEc2TagsTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := CreateEc2Tags(msvc, tt.resources, tt.tags); (err != nil) != tt.wantErr {
				t.Errorf("CreateEc2Tags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormEc2Tags(t *testing.T) {
	for n := range formEc2TagsTestCases {
		tt := formEc2TagsTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := FormEc2Tags(tt.tags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FormEc2Tags() = %v, want %v", got, tt.want)
			}
		})
	}
}
