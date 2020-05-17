package copyami

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mitchellh/mapstructure"
)

func (r *Resource) Prepare(rawConfig map[string]interface{}) error {
	var c Config

	var md mapstructure.Metadata

	if err := mapstructure.DecodeMetadata(rawConfig, &c, &md); err != nil {
		return err
	}

	r.Config = c
	r.Config.SrcAmiInfo = prepareSrcAmiInfo(r.Config.Source)

	return nil
}

func prepareSrcAmiInfo(rawSrcAmiInfo RawSrcAmiInfo) SrcAmiInfo {
	amiFilters := make([]*ec2.Filter, 0)

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
