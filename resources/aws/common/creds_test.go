package common

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
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
