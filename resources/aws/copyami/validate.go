package copyami

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/proffer/common/validator"
	"github.com/proffer/components"
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

	cs := validator.CustomStruct{Struct: c}
	if errs := validator.CheckRequiredFieldsInStruct(cs); len(errs) != 0 {
		clogger.Errorf("Missing/Empty key(s) found in the resource: [%s]", *r.Name)
		clogger.Fatal(errs)
	}

	r.validateConfigSource()
	r.validateConfigTarget()

	return nil
}

func (r *Resource) validateConfigSource() {
	sourceType := reflect.TypeOf(r.Config.Source)

	cs := validator.CustomStruct{Struct: r.Config.Source}
	if errs := validator.CheckRequiredFieldsInStruct(cs); len(errs) != 0 {
		clogger.Errorf("Missing/Empty key(s) found in the resource: [%s]", *r.Name)
		clogger.Fatal(errs)
	}

	if r.Config.Source.RoleArn != nil {
		sf, _ := sourceType.FieldByName("RoleArn")

		if !validator.IsAWSRoleARN(*r.Config.Source.RoleArn) {
			clogger.Fatalf("Invalid Role ARN [%s] passed in [%s] property of Resource: [%s]",
				*r.Config.Source.RoleArn, sf.Tag.Get("chain"), *r.Name)
		}
	}

	if !validator.IsAWSRegion(*r.Config.Source.Region) {
		sf, _ := sourceType.FieldByName("Region")
		clogger.Fatalf("Invalid AWS Region [%s] passed In [%s] property of Resource: [%s]",
			*r.Config.Source.Region, sf.Tag.Get("chain"), *r.Name)
	}

	sf, _ := sourceType.FieldByName("AmiFilters")

	for filterName, filterValue := range r.Config.Source.AmiFilters {
		if filterValue == nil {
			clogger.Fatalf("Missing value for AMI Filter [%s] in [%s] property of Resource: [%s]",
				*filterName, sf.Tag.Get("chain"), *r.Name)
		}

		switch *filterName {
		case "image-id":
			if !validator.IsAWSAMIID(*filterValue) {
				clogger.Fatalf("Invalid AWS AMI ID [%s] passed in [%s] property of Resource: [%s]",
					*filterValue, sf.Tag.Get("chain"), *r.Name)
			}
		case "name":
			if !validator.IsAWSAMIName(*filterValue) {
				clogger.Fatalf("Invalid AWS AMI Name [%s] passed in [%s] property of Resource: [%s]",
					*filterValue, sf.Tag.Get("chain"), *r.Name)
			}
		default:
			if strings.Contains(*filterName, "tag:") {
				if !validator.IsAWSTagKey(*filterName) {
					clogger.Fatalf("Invalid AWS Tag Key [%s] passed in [%s] property of Resource: [%s]",
						*filterName, sf.Tag.Get("chain"), *r.Name)
				}

				if !validator.IsAWSTagValue(*filterValue) {
					clogger.Fatalf("Invalid AWS Tag Value [%s] passed in [%s] property of Resource: [%s]",
						*filterValue, sf.Tag.Get("chain"), *r.Name)
				}
			}
		}
	}
}

func (r *Resource) validateConfigTarget() {
	targetType := reflect.TypeOf(r.Config.Target)

	cs := validator.CustomStruct{Struct: r.Config.Target}
	if errs := validator.CheckRequiredFieldsInStruct(cs); len(errs) != 0 {
		clogger.Errorf("Missing/Empty key(s) found in the resource: [%s]", *r.Name)
		clogger.Fatal(errs)
	}

	sf1, _ := targetType.FieldByName("Regions")

	for _, region := range r.Config.Target.Regions {
		if !validator.IsAWSRegion(*region) {
			clogger.Fatalf("Invalid AWS Region [%s] passed in [%s] property of Resource: [%s]",
				*region, sf1.Tag.Get("chain"), *r.Name)
		}
	}

	sf2, _ := targetType.FieldByName("AddExtraTags")

	for tagKey, tagValue := range r.Config.Target.AddExtraTags {
		if tagValue == nil {
			continue
		}

		if !validator.IsAWSTagKey(*tagKey) {
			clogger.Fatalf("Invalid AWS Tag Key [%s] passed in [%s] property of Resource: [%s]",
				*tagKey, sf2.Tag.Get("chain"), *r.Name)
		}

		if !validator.IsAWSTagValue(*tagValue) {
			clogger.Fatalf("Invalid AWS Tag Value [%s] passed in [%s] property of Resource: [%s]",
				*tagValue, sf2.Tag.Get("chain"), *r.Name)
		}
	}
}
