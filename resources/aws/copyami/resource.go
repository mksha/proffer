package copyami

import (
	"log"
	"os"

	clog "example.com/proffer/common/clogger"
	awscommon "example.com/proffer/resources/aws/common"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mitchellh/mapstructure"
)

var (
	clogger = clog.New(os.Stdout, "aws-copyami | ", log.Lmsgprefix)
)

type RawSrcAmiInfo struct {
	Profile    *string             `yaml:"profile"`
	RoleArn    *string             `yaml:"roleArn"`
	Region     *string             `yaml:"region"`
	AmiFilters map[*string]*string `yaml:"amiFilters"`
}

type SrcAmiInfo struct {
	CredsInfo map[string]string
	AccountID *string
	Region    *string
	Filters   []*ec2.Filter
	Image     *ec2.Image
}

type Target struct {
	Regions               []*string           `yaml:"regions"`
	CopyTagsAcrossRegions bool                `yaml:"copyTagsAcrossRegions"`
	AddExtraTags          map[*string]*string `yaml:"addExtraTags"`
}

type Config struct {
	Source     RawSrcAmiInfo `yaml:"source"`
	Target     Target        `yaml:"target"`
	SrcAmiInfo SrcAmiInfo    `yaml:"-"`
}

type Resource struct {
	Config Config `yaml:"config"`
}

func (r *Resource) Prepare(rawConfig map[string]interface{}) error {
	var c Config

	if err := mapstructure.Decode(rawConfig, &c); err != nil {
		return err
	}

	r.Config = c
	r.Config.SrcAmiInfo = prepareSrcAmiInfo(r.Config.Source)

	return nil
}

func (r *Resource) Run() error {
	target := r.Config.Target

	targetInfo := TargetInfo{
		Regions:  target.Regions,
		CopyTags: target.CopyTagsAcrossRegions,
		Tags:     awscommon.FormEc2Tags(target.AddExtraTags),
	}

	if err := apply(r.Config.SrcAmiInfo, targetInfo); err != nil {
		return err
	}

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
