package shareami

import (
	"reflect"
	"strconv"
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

	if errs := validator.CheckRequiredFieldsInStruct(c); len(errs) != 0 {
		clogger.Errorf("Missing/Empty key(s) found in the resource: [%s]", *r.Name)
		clogger.Fatal(errs)
	}

	r.validateConfigSource()
	r.validateConfigTarget()

	return nil
}

func (r *Resource) validateConfigSource() {
	sourceType := reflect.TypeOf(r.Config.Source)

	if errs := validator.CheckRequiredFieldsInStruct(r.Config.Source); len(errs) != 0 {
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
	sf, _ := targetType.FieldByName("AccountRegionMappingList")

	if errs := validator.CheckRequiredFieldsInStruct(r.Config.Target); len(errs) != 0 {
		clogger.Errorf("Missing/Empty key(s) found in the resource: [%s]", *r.Name)
		clogger.Fatal(errs)
	}

	for index, accountRegionMapping := range r.Config.Target.AccountRegionMappingList {
		acrMappingType := reflect.TypeOf(accountRegionMapping)

		if errs := validator.CheckRequiredFieldsInStruct(accountRegionMapping, index); len(errs) != 0 {
			clogger.Errorf("Missing/Empty key(s) found in the resource: [%s]", *r.Name)
			clogger.Fatal(errs)
		}

		if accountRegionMapping.Profile == nil && accountRegionMapping.RoleArn == nil {
			clogger.Fatalf("Need to specify one of the keys [profile or roleArn] in [%s.%v] item of Resource: [%s]",
				sf.Tag.Get("chain"), index, *r.Name)
		}

		if !validator.IsAWSAccountID(strconv.Itoa(accountRegionMapping.AccountID)) {
			sf, _ := acrMappingType.FieldByName("AccountID")
			chain := strings.Replace(sf.Tag.Get("chain"), ".N.", "."+strconv.Itoa(index)+".", 1)
			clogger.Fatalf("Invalid AWS Account ID [%v] passed in [%s] property of Resource: [%s]",
				accountRegionMapping.AccountID, chain, *r.Name)
		}

		sf1, _ := acrMappingType.FieldByName("Regions")
		chain := strings.Replace(sf1.Tag.Get("chain"), ".N.", "."+strconv.Itoa(index)+".", 1)

		for _, region := range accountRegionMapping.Regions {
			if !validator.IsAWSRegion(*region) {
				clogger.Fatalf("Invalid AWS Region [%s] passed in [%s] property of Resource: [%s]",
					*region, chain, *r.Name)
			}
		}

		sf2, _ := acrMappingType.FieldByName("AddExtraTags")
		chain = strings.Replace(sf2.Tag.Get("chain"), ".N.", "."+strconv.Itoa(index)+".", 1)

		for tagKey, tagValue := range accountRegionMapping.AddExtraTags {
			if tagValue == nil {
				continue
			}

			if !validator.IsAWSTagKey(*tagKey) {
				clogger.Fatalf("Invalid AWS Tag Key [%s] passed in [%s] property of Resource: [%s]",
					*tagKey, chain, *r.Name)
			}

			if !validator.IsAWSTagValue(*tagValue) {
				clogger.Fatalf("Invalid AWS Tag Value [%s] passed in [%s] property of Resource: [%s]",
					*tagValue, chain, *r.Name)
			}
		}
	}

	sf, _ = targetType.FieldByName("CommonRegions")

	for _, region := range r.Config.Target.CommonRegions {
		if !validator.IsAWSRegion(*region) {
			clogger.Fatalf("Invalid AWS Region [%s] passed in [%s] property of Resource: [%s]",
				*region, sf.Tag.Get("chain"), *r.Name)
		}
	}
}
