package copyami

import (
	"log"
	"os"

	clog "example.com/proffer/common/clogger"
	awscommon "example.com/proffer/resources/aws/common"
	"github.com/mitchellh/mapstructure"
)

var (
	clogger = clog.New(os.Stdout, "aws-copyami | ", log.Lmsgprefix)
)

type Target struct {
	Regions               []*string           `yaml:"regions"`
	CopyTagsAcrossRegions bool                `yaml:"copyTagsAcrossRegions"`
	AddExtraTags          map[*string]*string `yaml:"addExtraTags"`
}

type Config struct {
	Source     awscommon.RawSrcAmiInfo `yaml:"source"`
	Target     Target                  `yaml:"target"`
	SrcAmiInfo awscommon.SrcAmiInfo    `yaml:"-"`
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
	r.Config.SrcAmiInfo = awscommon.PrepareSrcAmiInfo(r.Config.Source)

	return nil
}

func (r *Resource) Run() error {
	target := r.Config.Target

	targetInfo := TargetInfo{
		Regions:  target.Regions,
		CopyTags: target.CopyTagsAcrossRegions,
		Tags:     awscommon.FormEc2Tags(target.AddExtraTags),
	}

	if err := copyAmi(r.Config.SrcAmiInfo, targetInfo); err != nil {
		return err
	}

	return nil
}
