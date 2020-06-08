package common

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts"
)

func TestIsCredsExpired(t *testing.T) {
	tests := []struct {
		name string
		msvc *mockSTSClient
		want bool
	}{
		{
			name: "valid creds",
			msvc: &mockSTSClient{
				Error: nil,
			},
			want: false,
		},
		{
			name: "expired creds",
			msvc: &mockSTSClient{
				Error: awserr.New("ExpiredToken", "Creds have expired", errors.New("ExpiredToken")),
			},
			want: true,
		},
	}

	for n := range tests {
		tt := tests[n]
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCredsExpired(tt.msvc); got != tt.want {
				t.Errorf("IsCredsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCallerInfo(t *testing.T) {
	tests := []struct {
		name    string
		svc     *mockSTSClient
		want    *sts.GetCallerIdentityOutput
		wantErr bool
	}{
		{
			name: "valid account request",
			svc: &mockSTSClient{
				Error: nil,
			},
			want: &sts.GetCallerIdentityOutput{
				Account: aws.String("123456789012"),
				Arn:     aws.String("arn:aws::iam:123456789012:role/test"),
				UserId:  aws.String("123@testuser"),
			},
			wantErr: false,
		},
		{
			name: "expired creds",
			svc: &mockSTSClient{
				Error: awserr.New("ExpiredToken", "Creds have expired", errors.New("ExpiredToken")),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for n := range tests {
		tt := tests[n]
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCallerInfo(tt.svc)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCallerInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCallerInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAccountAlias(t *testing.T) {
	tests := []struct {
		name    string
		svc     iamiface.IAMAPI
		want    *string
		wantErr bool
	}{
		{
			name: "valid account request",
			svc: &mockIAMClient{
				Error: nil,
			},
			want:    aws.String("test-account"),
			wantErr: false,
		},
		{
			name: "expired creds",
			svc: &mockIAMClient{
				Error: awserr.New("ExpiredToken", "Creds have expired", errors.New("ExpiredToken")),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for n := range tests {
		tt := tests[n]
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAccountAlias(tt.svc)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAccountAlias() = %v, want %v", *got, *tt.want)
			}
		})
	}
}
