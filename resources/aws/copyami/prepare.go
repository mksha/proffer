package copyami

import (
	"github.com/mitchellh/mapstructure"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (r *Resource) Prepare(rawConfig map[string]interface{}) error {
	var c Config

	if err := mapstructure.Decode(rawConfig, &c); err != nil {
		return err
	}

	if c.Other != nil {
		clogger.Fatalf("Please remove invalid sections/properties: %s", c.Other)
	}

	r.Config = c
	r.Config.SrcAmiInfo = prepareSrcAmiInfo(r.Config.Source)

	return nil
}

func prepareSrcAmiInfo(rawSrcAmiInfo RawSrcAmiInfo) SrcAmiInfo {
	var amiFilters []*ec2.Filter

	for filterName, filterValue := range rawSrcAmiInfo.AmiFilters {
		f := &ec2.Filter{
			Name:   filterName,
			Values: []*string{filterValue},
		}
		amiFilters = append(amiFilters, f)
	}

	srcAmiInfo := SrcAmiInfo{
		Region:    rawSrcAmiInfo.Region,
		Filters:   amiFilters,
		CredsInfo: make(map[string]string, 2),
	}

	if rawSrcAmiInfo.RoleArn != nil {
		srcAmiInfo.CredsInfo["getCredsUsing"] = "roleArn"
		srcAmiInfo.CredsInfo["roleArn"] = *rawSrcAmiInfo.RoleArn
	} else if rawSrcAmiInfo.Profile != nil {
		srcAmiInfo.CredsInfo["getCredsUsing"] = "profile"
		srcAmiInfo.CredsInfo["profile"] = *rawSrcAmiInfo.Profile
	}

	return srcAmiInfo
}
