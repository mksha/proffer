package common

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// test cases for IsAmiExist function.
var isAmiExistTestCases = []struct {
	name    string
	filters []*ec2.Filter
	want    bool
	wantErr bool
}{
	{
		name:    "call with no filters",
		filters: []*ec2.Filter{},
		want:    false,
		wantErr: true,
	},
	{
		name: "call with valid filters",
		filters: []*ec2.Filter{
			{
				Name: aws.String("name"),
				Values: []*string{
					aws.String("test-image"),
				},
			},
			{
				Name: aws.String("tag:Purpose"),
				Values: []*string{
					aws.String("testing"),
				},
			},
		},
		want:    true,
		wantErr: false,
	},
	{
		name: "filters that will not give any ami (give empty image list)",
		filters: []*ec2.Filter{
			{
				Name: aws.String("unknown"),
				Values: []*string{
					aws.String("unknown"),
				},
			},
		},
		want:    false,
		wantErr: false,
	},
}

// test cases for CreateEc2Tags function.
var createEc2TagsTestCases = []struct {
	name      string
	resources []*string
	tags      []*ec2.Tag
	wantErr   bool
}{
	{
		name:      "empty resources and empty tags list",
		resources: []*string{},
		tags:      []*ec2.Tag{},
		wantErr:   true,
	},
	{
		name:      "empty resources list",
		resources: []*string{},
		tags: []*ec2.Tag{
			{
				Key:   aws.String("test"),
				Value: aws.String("test"),
			},
		},
		wantErr: true,
	},
	{
		name:      "empty tags list",
		resources: []*string{aws.String("ami-123456789012")},
		tags:      []*ec2.Tag{},
		wantErr:   true,
	},
	{
		name:      "Valid resources and tags list",
		resources: []*string{aws.String("ami-123456789012")},
		tags: []*ec2.Tag{
			{
				Key:   aws.String("test"),
				Value: aws.String("test"),
			},
		},
		wantErr: false,
	},
}

// test cases for FormEc2Tags function.
var formEc2TagsTestCases = []struct {
	name string
	tags map[*string]*string
	want []*ec2.Tag
}{
	{
		name: "empty map of tags",
		tags: map[*string]*string{},
		want: []*ec2.Tag{},
	},
	{
		name: "valid map of tags",
		tags: map[*string]*string{
			aws.String("testkey"):    aws.String("testvalue"),
			aws.String("test key2 "): aws.String("test value 12187^B*&@*$*@%$*O$"),
		},
		want: []*ec2.Tag{
			{
				Key:   aws.String("testkey"),
				Value: aws.String("testvalue"),
			},
			{
				Key:   aws.String("test key2 "),
				Value: aws.String("test value 12187^B*&@*$*@%$*O$"),
			},
		},
	},
}

// test cases for IsError function.
var isErrorTestCases = []struct {
	name    string
	err     error
	want    bool
	wantErr bool
}{
	{
		name:    "no error",
		err:     nil,
		want:    false,
		wantErr: false,
	},
	{
		name:    "general error",
		err:     errors.New("general error"),
		want:    true,
		wantErr: true,
	},
	{
		name:    "aws type error",
		err:     awserr.New("test-code", "test-msg", errors.New("test error")),
		want:    true,
		wantErr: true,
	},
	{
		name:    "aws RequestExpired type error",
		err:     awserr.New("RequestExpired", "test-msg", errors.New("test error")),
		want:    true,
		wantErr: true,
	},
}

//test cases for GetAmiInfo function.
var getAmiInfoTestCases = []struct {
	name    string
	filters []*ec2.Filter
	want    []*ec2.Image
	wantErr bool
}{
	{
		name:    "empty filters",
		filters: []*ec2.Filter{},
		want:    nil,
		wantErr: true,
	},
	{
		name: "filters that will not give any ami (give empty image list)",
		filters: []*ec2.Filter{
			{
				Name: aws.String("unknown"),
				Values: []*string{
					aws.String("unknown"),
				},
			},
		},
		want:    nil,
		wantErr: true,
	},
	{
		name: "call with valid filters",
		filters: []*ec2.Filter{
			{
				Name: aws.String("name"),
				Values: []*string{
					aws.String("test-image"),
				},
			},
			{
				Name: aws.String("tag:Purpose"),
				Values: []*string{
					aws.String("testing"),
				},
			},
		},
		want: []*ec2.Image{{
			ImageId: aws.String("test-image-id"),
		}},
		wantErr: false,
	},
}
