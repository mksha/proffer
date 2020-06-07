package common

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	clog "github.com/proffer/common/clogger"
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

// GetAwsSessWithAssumeRole returns a pointer to aws session which got the credentials by assuming the given role arn.
// It also returns an error if there was any.
func GetAwsSessWithAssumeRole(roleArn string) (*session.Session, error) {
	sessPtr, err := GetAwsSessWithDefaultCreds()

	if err != nil {
		return nil, err
	}

	svc := sts.New(sessPtr)

	callerInfo, err := GetCallerInfo(svc)
	if err != nil {
		return nil, err
	}

	clogger.Debugf("Will Assume IAM Role Using Creds Of Identity: %s", *callerInfo.Arn)

	creds := stscreds.NewCredentials(sessPtr, roleArn)
	config := aws.Config{Credentials: creds}
	newSessPtr := createSession(session.Options{Config: config})

	if _, err := newSessPtr.Config.Credentials.Get(); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, fmt.Errorf("%s: %s", aerr.Code(), aerr.Message())
		}
	}

	if ok := IsCredsExpired(svc); ok {
		return nil, fmt.Errorf("AwsAssumeRoleCredsExpired: Provided AWS Assume Role's Credentials Have Expired")
	}

	return newSessPtr, nil
}

// GetAwsSessWithProfile returns a pointer to profile that has got creds from aws profile.
// It also returns an error if there was any.
func GetAwsSessWithProfile(profile string) (*session.Session, error) {
	sessPtr := createSession(session.Options{Profile: profile})

	if _, err := sessPtr.Config.Credentials.Get(); err != nil {
		return nil, fmt.Errorf("AWSProfileDoesNotExist: Failed To Retrieve Credentials From AWS Profile '%s'", profile)
	}

	svc := sts.New(sessPtr)
	if ok := IsCredsExpired(svc); ok {
		return nil, fmt.Errorf("AwsProfileCredsExpired: AWS Profile '%s's Credentials Have Expired", profile)
	}

	return sessPtr, nil
}

// GetAwsSessWithDefaultCreds returns a pointer to session which has got aws creds from default methods.
// It also returns an error if there was any.
func GetAwsSessWithDefaultCreds() (*session.Session, error) {
	sessPtr := createSession(session.Options{SharedConfigState: session.SharedConfigEnable})

	if _, err := sessPtr.Config.Credentials.Get(); err != nil {
		return nil, fmt.Errorf("NoDefaultCredProviderExist: No Default Credential Provider Exists")
	}

	svc := sts.New(sessPtr)
	if ok := IsCredsExpired(svc); ok {
		return nil, fmt.Errorf("CredsExpired: Default AWS Provider's Credentials Have Expired")
	}

	return sessPtr, nil
}

// IsCredsExpired returns true if given session is having expired credentials.
func IsCredsExpired(svc stsiface.STSAPI) bool {
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

// GetCallerInfo returns the caller identity from the given service client session.
// It also returns an error if there was any.
func GetCallerInfo(svc stsiface.STSAPI) (*sts.GetCallerIdentityOutput, error) {
	result, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAccountAlias returns the AWS Account Alias for the given service client session.
// It also returns an error if there was any.
func GetAccountAlias(svc iamiface.IAMAPI) (*string, error) {
	input := &iam.ListAccountAliasesInput{}
	result, err := svc.ListAccountAliases(input)

	if ok, err := IsError(err); ok {
		return nil, err
	}

	return result.AccountAliases[0], nil
}

// GetAwsSession returns a session pointer based on given credential provider.
// It also returns an error if there was any.
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
