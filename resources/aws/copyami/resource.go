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
	Profile    *string             `mapstructure:"profile" required:"false"`
	RoleArn    *string             `mapstructure:"roleArn" required:"false"`
	Region     *string             `mapstructure:"region" required:"true" chain:"config.source.region"`
	AmiFilters map[*string]*string `mapstructure:"amiFilters" required:"true" chain:"config.source.amiFilters"`
}

type SrcAmiInfo struct {
	CredsInfo map[string]string
	AccountID *string
	Region    *string
	Filters   []*ec2.Filter
	Image     *ec2.Image
}

type Target struct {
	Regions               []*string           `mapstructure:"regions" required:"true" chain:"config.target.regions"`
	CopyTagsAcrossRegions bool                `mapstructure:"copyTagsAcrossRegions"`
	AddExtraTags          map[*string]*string `mapstructure:"addExtraTags"`
}

type Config struct {
	Source     RawSrcAmiInfo `mapstructure:"source" required:"true" chain:"config.source"`
	Target     Target        `mapstructure:"target" required:"true" chain:"config.target"`
	SrcAmiInfo SrcAmiInfo    `mapstructure:"-"`
}

type Resource struct {
	Name   *string `required:"true"`
	Type   *string `required:"true"`
	Config Config  `mapstructure:"config" required:"true"`
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
