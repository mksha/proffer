package common

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

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
