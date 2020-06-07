package copyami

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/ec2"
	clog "github.com/proffer/common/clogger"
	awscommon "github.com/proffer/resources/aws/common"
	"gopkg.in/yaml.v2"
)

var (
	clogger = clog.New(os.Stdout, "aws-copyami | ", log.Lmsgprefix)
)

// RawSrcAmiInfo represents the raw source ami information.
type RawSrcAmiInfo struct {
	Profile    *string             `mapstructure:"profile" required:"false" chain:"config.source.profile"`
	RoleArn    *string             `mapstructure:"roleArn" required:"false" chain:"config.source.roleArn"`
	Region     *string             `mapstructure:"region" required:"true" chain:"config.source.region"`
	AmiFilters map[*string]*string `mapstructure:"amiFilters" required:"true" chain:"config.source.amiFilters"`
}

// SrcAmiInfo represents the parsed source ami information.
type SrcAmiInfo struct {
	CredsInfo map[string]string
	AccountID *string
	Region    *string
	Filters   []*ec2.Filter
	Image     *ec2.Image
}

// Target represents the target configuration needed by aws-copyami resource.
type Target struct {
	Regions               []*string           `mapstructure:"regions" required:"true" chain:"config.target.regions"`
	CopyTagsAcrossRegions bool                `mapstructure:"copyTagsAcrossRegions" chain:"config.target.copyTagsAcrossRegions"`
	AddExtraTags          map[*string]*string `mapstructure:"addExtraTags" chain:"config.target.addExtraTags"`
}

// Config represents the configuration needed by aws-copyami resource in go struct.
type Config struct {
	Source     RawSrcAmiInfo `mapstructure:"source" required:"true" chain:"config.source"`
	Target     Target        `mapstructure:"target" required:"true" chain:"config.target"`
	SrcAmiInfo SrcAmiInfo    `mapstructure:"-"`
}

// SrcImage represent the source ami information used for inventory generation.
type SrcImage struct {
	ID     *string `yaml:"id"`
	Name   *string `yaml:"name"`
	Region *string `yaml:"region"`
}

// AmiMeta represent the AWS AMI metadata used for inventory generation.
type AmiMeta struct {
	ID   *string `yaml:"id"`
	Name *string `yaml:"name"`
}

// AccountMeta represents the AWS Account metadata like id, alias, etc.
type AccountMeta struct {
	ID    *string
	Alias *string
}

// Record represents the inventory record for aws-copyami resource.
type Record struct {
	AccountMeta  AccountMeta         `yaml:"account"`
	SrcImage     SrcImage            `yaml:"sourceImage"`
	TargetImages map[*string]AmiMeta `yaml:"targetImages"`
}

// Resource represents the aws-copyami resource.
type Resource struct {
	Name   *string `required:"true"`
	Type   *string `required:"true"`
	Config Config  `mapstructure:"config" required:"true"`
	Record Record  `mapstructure:"-"`
}

// Run applies the resource specific configuration.
func (r *Resource) Run() error {
	target := r.Config.Target

	targetInfo := TargetInfo{
		Regions:  target.Regions,
		CopyTags: target.CopyTagsAcrossRegions,
		Tags:     awscommon.FormEc2Tags(target.AddExtraTags),
	}

	if err := r.apply(r.Config.SrcAmiInfo, targetInfo); err != nil {
		return err
	}

	clogger.Info(r.Record)
	bs, _ := yaml.Marshal(r.Record)
	clogger.Info(string(bs))

	file, err := os.Create("inventory.yml")
	if err != nil {
		return err
	}

	_, err = file.Write(bs)
	if err != nil {
		return err
	}

	return nil
}
