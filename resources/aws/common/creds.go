package common

import (
	// "fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// createSession is used to create aws session with given configuration.
func createSession(sessOpts session.Options) (sessPtr *session.Session) {
	// Create a Session with a custom config
	sessPtr = session.Must(session.NewSessionWithOptions(sessOpts))
	return
}

func GetAwsSessWithProfile(profile, region string) (sessPtr *session.Session) {
	sessOpts := session.Options{
		// Specify profile to load for the session's config
		Profile: profile,
		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region: aws.String(region),
		},
		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	}
	sessPtr = createSession(sessOpts)
	return
}

func GetAwsSessWithDefaultCreds() (sessPtr *session.Session) {
	creds := credentials.NewEnvCredentials()

	if _, err := creds.Get(); err != nil {
		log.Fatalln(" AWSEnvVarsDoesNotExist: Failed To Retrive Credentials From Env Vars")
	}

	if ok := creds.IsExpired(); ok {
		log.Fatalln(" AWS Environment Credentials have Expired")
	}

	sessPtr = createSession(session.Options{})
	return
}

// func validateCreds(sess *session.Session) error {
// 	svc := sts.New(sess)
// 	input := &sts.GetCallerIdentityInput{}
// 	_, err := svc.GetCallerIdentity(input)

// 	if err != nil {
// 		if aerr, ok := err.(awserr.Error); ok {
// 			if aerr.Code() == "ExpiredToken" {
// 				return fmt.Errorf("%s: %s", aerr.Code(), aerr.Message())
// 			}
// 		}

// 		return err
// 	}

// 	return nil
// }
