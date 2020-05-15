package copyami

import (
	"fmt"
	// "reflect"

	"example.com/proffer/common/validator"
	"example.com/proffer/components"
	"github.com/mitchellh/mapstructure"
)

func (r *Resource) Validate(rawResource components.RawResource) error {
	var c Config

	r.Name = &rawResource.Name
	r.Type = &rawResource.Type

	var md mapstructure.Metadata

	if err := mapstructure.DecodeMetadata(rawResource.Config, &c, &md); err != nil {
		return err
	}

	r.Config = c

	if len(md.Unused) != 0 {
		clogger.Errorf("Invalid key(s) in config section of Resource: [%s]", *r.Name)
		clogger.Error(md.Unused)

		return fmt.Errorf("validation failed")
	}

	if errs := validator.CheckRequiredFieldsInStruct(c); len(errs) != 0 {
		clogger.Errorf("Missing/Empty key(s) found in the resource: [%s]", *r.Name)
		clogger.Fatal(errs)
	}

	r.validateConfigSource()
	r.validateConfigTarget()

	return nil
}

func (r *Resource) validateConfigSource() {
	if errs := validator.CheckRequiredFieldsInStruct(r.Config.Source); len(errs) != 0 {
		clogger.Errorf("Missing/Empty key(s) found in the resource: [%s]", *r.Name)
		clogger.Fatal(errs)
	}

	if r.Config.Source.RoleArn != nil {
		if !validator.IsAWSRoleArn(*r.Config.Source.RoleArn) {
			clogger.Fatalf("Invalid Role ARN Passed In 'config.source.roleArn' Property Of Resource: [%s]", *r.Name)
		}
	}
}

func (r *Resource) validateConfigTarget() {
	if errs := validator.CheckRequiredFieldsInStruct(r.Config.Target); len(errs) != 0 {
		clogger.Errorf("Missing/Empty key(s) found in the resource: [%s]", *r.Name)
		clogger.Fatal(errs)
	}
}
