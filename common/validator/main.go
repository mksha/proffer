package validator

import (
	"fmt"
	"reflect"
	"regexp"
)

var (
	re = map[string]*regexp.Regexp{
		"awsRoleARN": regexp.MustCompile(`(?m)^arn:(aws(|-cn|-us-gov))?:iam::\d{12}:role/?[a-zA-Z_0-9+=,.@\-_/]+`),
	}
)

func IsZero(i interface{}) bool {
	v := reflect.ValueOf(i)
	return v.IsZero()
}

func CheckRequiredFieldsInStruct(i interface{}) []error {
	errs := make([]error, 0)
	v := reflect.ValueOf(i)

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
					err = fmt.Errorf(vt.Field(i).Tag.Get("chain"))
				} else if vt.Field(i).Tag.Get("metastructure") != "" {
					err = fmt.Errorf(vt.Field(i).Tag.Get("metastructure"))
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

func IsAWSRoleArn(roleARN string) bool {
	return len(re["awsRoleARN"].FindStringIndex(roleARN)) > 0
}

func isAWSRegion() bool {
	return true
}
