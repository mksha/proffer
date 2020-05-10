package common

import (
	// "fmt"
	"fmt"
	"log"
	"os"
	// "time"

	clog "example.com/proffer/common/clogger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	// "github.com/aws/aws-sdk-go/aws/credentials"
	// "github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	// "github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

var (
	clogger = clog.New(os.Stdout, "common | ", log.Lmsgprefix)
)

// createSession is used to create aws session with given configuration.
func createSession(sessOpts session.Options) (sessPtr *session.Session) {
	// Create a Session with a custom config
	sessPtr = session.Must(session.NewSessionWithOptions(sessOpts))
	return
}

func GetAwsSessWithAssumeRole(roleArn string) (*session.Session, error) {
	sessPtr, err := GetAwsSessWithDefaultCreds()

	if err != nil {
		return nil, err
	}

	callerAccountInfo, err := GetAccountInfo(sessPtr)
	if err != nil {
		return nil, err
	}

	clogger.Debugf("Will Assume IAM Role Using Creds Of Identity: %s", *callerAccountInfo.Arn)

	creds := stscreds.NewCredentials(sessPtr, roleArn)
	config := aws.Config{Credentials: creds}
	newSessPtr := createSession(session.Options{Config: config})

	if _, err := newSessPtr.Config.Credentials.Get(); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, fmt.Errorf("%s: %s", aerr.Code(), aerr.Message())
		}
	}

	if ok := IsCredsExpired(newSessPtr); ok {
		return nil, fmt.Errorf("AwsAssumeRoleCredsExpired: Provided AWS Assume Role's Credentials Have Expired")
	}

	return newSessPtr, nil
}

func GetAwsSessWithProfile(profile string) (*session.Session, error) {
	sessPtr := createSession(session.Options{Profile: profile})

	if _, err := sessPtr.Config.Credentials.Get(); err != nil {
		return nil, fmt.Errorf("AWSProfileDoesNotExist: Failed To Retrieve Credentials From AWS Profile '%s'", profile)
	}

	if ok := IsCredsExpired(sessPtr); ok {
		return nil, fmt.Errorf("AwsProfileCredsExpired: AWS Profile '%s's Credentials Have Expired", profile)
	}

	return sessPtr, nil
}

func GetAwsSessWithDefaultCreds() (*session.Session, error) {
	sessPtr := createSession(session.Options{SharedConfigState: session.SharedConfigEnable})

	if _, err := sessPtr.Config.Credentials.Get(); err != nil {
		return nil, fmt.Errorf("NoDefaultCredProviderExist: No Default Credential Provider Exists")
	}

	if ok := IsCredsExpired(sessPtr); ok {
		return nil, fmt.Errorf("CredsExpired: Default AWS Provider's Credentials Have Expired")
	}

	return sessPtr, nil
}

func IsCredsExpired(sessPtr *session.Session) bool {
	svc := sts.New(sessPtr)
	_, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "ExpiredToken" {
				return true
			}
		}

		clogger.Fatal(err)
	}

	return false
}

func GetAccountInfo(sessPtr *session.Session) (*sts.GetCallerIdentityOutput, error) {
	svc := sts.New(sessPtr)
	result, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetAwsSession(credProviderInfo map[string]string) (*session.Session, error) {
	switch credProviderInfo["getCredsUsing"] {
	case "profile":
		clogger.Debugf("Will Get Creds Using AWS Profile: '%s'", credProviderInfo["profile"])
		return GetAwsSessWithProfile(credProviderInfo["profile"])
	case "roleArn":
		clogger.Debugf("Will Get Creds By Assuming AWS IAM Role: '%s'", credProviderInfo["roleArn"])
		return GetAwsSessWithAssumeRole(credProviderInfo["roleArn"])
	default:
		credLoadOrder := `
	 * Environment Variables
	 * Shared Credentials file
	 * Shared Configuration file
	 * EC2 Instance Metadata (credentials only)`
		clogger.Debugf("Will Attempt To Load Creds Using One Of the Providers Found In Following Order: %s", credLoadOrder)

		return GetAwsSessWithDefaultCreds()
	}
}
