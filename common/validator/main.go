package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	re = map[string]*regexp.Regexp{
		"awsRoleARN":   regexp.MustCompile(`^arn:(aws(|-cn|-us-gov))?:iam::\d{12}:role/[a-zA-Z_0-9+=,.@\-_/]{1,64}$`),
		"awsRegion":    regexp.MustCompile(`(^(us(-gov)?)-(east|west)-(1|2)$)|(^eu-(central-1$|west-(1|2|3)$|(south|north)-1$))|(^cn-(north(west)?-1$))|(^ap-(east-1$|south(east)?-(1|2)$|northeast-(1|2|3)$))|(^ca-central-1$)|(^(af|me)-south-1$)|(^sa-east-1$)`), //nolint:lll
		"awsAMIName":   regexp.MustCompile(`^[a-zA-Z0-9\-_.\(\)\/]{3,128}$`),
		"awsAMIID":     regexp.MustCompile(`^ami-[0-9a-z]{8,17}$`),
		"awsTagKey":    regexp.MustCompile(`^.{1,127}$`),
		"awsTagValue":  regexp.MustCompile(`^.{0,127}$`),
		"awsAccountID": regexp.MustCompile(`^\d{12}$`),
	}
)

// IsZero checks if the given value is zero values of its type.
func IsZero(i interface{}) bool {
	v := reflect.ValueOf(i)
	return v.IsZero()
}

//nolint:nestif,gocritic
// CheckRequiredFieldsInStruct checks if the given structs's required field(s) are set or not.
func CheckRequiredFieldsInStruct(customStruct interface{}, index ...int) []error {
	errs := make([]error, 0)
	v := reflect.ValueOf(customStruct)

	if v.Kind() == reflect.Struct {
		vt := v.Type() // to get the filed info like tags, pkgpath etc

		for i := 0; i < v.NumField(); i++ {
			if vt.Field(i).PkgPath != "" {
				continue // Private field
			}

			// check if the required tag is true
			if vt.Field(i).Tag.Get("required") == "false" || vt.Field(i).Tag.Get("required") == "" {
				continue // optional field
			}

			if v.Field(i).IsZero() {
				var err error

				if vt.Field(i).Tag.Get("chain") != "" {
					if len(index) != 0 {
						msg := strings.Replace(vt.Field(i).Tag.Get("chain"), ".N.", "."+strconv.Itoa(index[0])+".", 1)
						err = fmt.Errorf(msg)
					} else {
						err = fmt.Errorf(vt.Field(i).Tag.Get("chain"))
					}
				} else if vt.Field(i).Tag.Get("mapstructure") != "" {
					err = fmt.Errorf(vt.Field(i).Tag.Get("mapstructure"))
				} else {
					err = fmt.Errorf(vt.Field(i).Tag.Get("yaml"))
				}

				if err != nil {
					errs = append(errs, err)
				}
			}
		}
	}

	return errs
}

// IsAWSRoleARN validates if arn is a valid AWS Role ARN or not.
func IsAWSRoleARN(arn string) bool {
	return re["awsRoleARN"].MatchString(arn)
}

// IsAWSRegion validates if region is a valid AWS Region or not.
func IsAWSRegion(region string) bool {
	return re["awsRegion"].MatchString(region)
}

// IsAWSAMIID validates if id is a valid AWS AMI ID or not.
func IsAWSAMIID(id string) bool {
	return re["awsAMIID"].MatchString(id)
}

// IsAWSAMIName validates if name is a valid AWS AMI Name or not.
func IsAWSAMIName(name string) bool {
	return re["awsAMIName"].MatchString(name)
}

// IsAWSTagKey returns a boolean indicating whether the key is a valid AWS Tag Key or not.
func IsAWSTagKey(key string) bool {
	return re["awsTagKey"].MatchString(key)
}

// IsAWSTagValue returns a boolean indicating whether the value is a valid AWS Tag Value or not.
func IsAWSTagValue(value string) bool {
	return re["awsTagValue"].MatchString(value)
}

// IsAWSAccountID reports whether id is a valid AWS Account ID or not.
func IsAWSAccountID(id string) bool {
	return re["awsAccountID"].MatchString(id)
}
