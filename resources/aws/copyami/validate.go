package copyami

import (
	"fmt"
	"strings"

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
		if !validator.IsAWSRoleARN(*r.Config.Source.RoleArn) {
			clogger.Fatalf("Invalid Role ARN [%s] Passed In [config.source.roleArn] Property Of Resource: [%s]",
				*r.Config.Source.RoleArn, *r.Name)
		}
	}

	if !validator.IsAWSRegion(*r.Config.Source.Region) {
		clogger.Fatalf("Invalid AWS Region [%s] Passed In [config.source.region] Property Of Resource: [%s]",
			*r.Config.Source.Region, *r.Name)
	}

	for filterName, filterValue := range r.Config.Source.AmiFilters {
		if filterValue == nil {
			clogger.Fatalf("Missing Value For AMI Filter [%s] in [config.source.amiFilters] Property Of Resource: [%s]",
				*filterName, *r.Name)
		}

		switch *filterName {
		case "image-id":
			if !validator.IsAWSAMIID(*filterValue) {
				clogger.Fatalf("Invalid AWS AMI ID [%s] Passed In [config.source.amiFilters] Property Of Resource: [%s]",
					*filterValue, *r.Name)
			}
		case "name":
			if !validator.IsAWSAMIName(*filterValue) {
				clogger.Fatalf("Invalid AWS AMI Name [%s] Passed In [config.source.amiFilters] Property Of Resource: [%s]",
					*filterValue, *r.Name)
			}
		default:
			if strings.Contains(*filterName, "tag:") {
				if !validator.IsAWSTagKey(*filterName) {
					clogger.Fatalf("Invalid AWS Tag Key [%s] Passed In [config.source.amiFilters] Property Of Resource: [%s]",
						*filterName, *r.Name)
				}

				if !validator.IsAWSTagValue(*filterValue) {
					clogger.Fatalf("Invalid AWS Tag Value [%s] Passed In [config.source.amiFilters] Property Of Resource: [%s]",
						*filterValue, *r.Name)
				}
			}
		}
	}
}

func (r *Resource) validateConfigTarget() {
	if errs := validator.CheckRequiredFieldsInStruct(r.Config.Target); len(errs) != 0 {
		clogger.Errorf("Missing/Empty key(s) found in the resource: [%s]", *r.Name)
		clogger.Fatal(errs)
	}

	for _, region := range r.Config.Target.Regions {
		if !validator.IsAWSRegion(*region) {
			clogger.Fatalf("Invalid AWS Region [%s] Passed In [config.target.regions] Property Of Resource: [%s]",
				*region, *r.Name)
		}
	}

	for tagKey, tagValue := range r.Config.Target.AddExtraTags {
		if tagValue == nil {
			continue
		}

		if !validator.IsAWSTagKey(*tagKey) {
			clogger.Fatalf("Invalid AWS Tag Key [%s] Passed In [config.source.addExtraTags] Property Of Resource: [%s]",
				*tagKey, *r.Name)
		}

		if !validator.IsAWSTagValue(*tagValue) {
			clogger.Fatalf("Invalid AWS Tag Value [%s] Passed In [config.source.addExtraTags] Property Of Resource: [%s]",
				*tagValue, *r.Name)
		}
	}
}
