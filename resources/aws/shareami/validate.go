package shareami

import (
	"strconv"
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

	for index, accountRegionMapping := range r.Config.Target.AccountRegionMappingList {
		if errs := validator.CheckRequiredFieldsInStruct(accountRegionMapping, index); len(errs) != 0 {
			clogger.Errorf("Missing/Empty key(s) found in the resource: [%s]", *r.Name)
			clogger.Fatal(errs)
		}

		if accountRegionMapping.Profile == nil && accountRegionMapping.RoleArn == nil {
			clogger.Fatalf("Need To Specify One Of The Keys [profile or roleArn] In [config.target.accountRegionMappingList.%v] List Of Resource: [%s]",
				index, *r.Name)
		}

		if !validator.IsAWSAccountID(strconv.Itoa(accountRegionMapping.AccountID)) {
			clogger.Fatalf("Invalid AWS Account ID [%v] Passed In [config.target.accountRegionMappingList.%v.accountId] Property Of Resource: [%s]",
				accountRegionMapping.AccountID, index, *r.Name)
		}

		for _, region := range accountRegionMapping.Regions {
			if !validator.IsAWSRegion(*region) {
				clogger.Fatalf("Invalid AWS Region [%s] Passed In [config.target.accountRegionMappingList.%v.regions] Property Of Resource: [%s]",
					*region, index, *r.Name)
			}
		}

		for tagKey, tagValue := range accountRegionMapping.AddExtraTags {
			if tagValue == nil {
				continue
			}

			if !validator.IsAWSTagKey(*tagKey) {
				clogger.Fatalf("Invalid AWS Tag Key [%s] Passed In [config.target.accountRegionMappingList.%v.addExtraTags] Property Of Resource: [%s]",
					*tagKey, index, *r.Name)
			}

			if !validator.IsAWSTagValue(*tagValue) {
				clogger.Fatalf("Invalid AWS Tag Value [%s] Passed In [config.target.accountRegionMappingList.%v.addExtraTags] Property Of Resource: [%s]",
					*tagValue, index, *r.Name)
			}
		}
	}

	for _, region := range r.Config.Target.CommonRegions {
		if !validator.IsAWSRegion(*region) {
			clogger.Fatalf("Invalid AWS Region [%s] Passed In [config.target.commonRegions] Property Of Resource: [%s]",
				*region, *r.Name)
		}
	}
}
