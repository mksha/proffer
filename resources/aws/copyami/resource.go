package copyami

import (
	"log"
	"os"

	clog "example.com/proffer/common/clogger"
	awscommon "example.com/proffer/resources/aws/common"
	"github.com/aws/aws-sdk-go/service/ec2"
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
	Other  map[string]interface{} `mapstructure:",remain"`
}

type Resource struct {
	Config Config `yaml:"config"`
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
