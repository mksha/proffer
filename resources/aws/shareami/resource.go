package shareami

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/ec2"
	clog "github.com/proffer/common/clogger"
)

var (
	clogger = clog.New(os.Stdout, "aws-shareami | ", log.Lmsgprefix)
)

type RawSrcAmiInfo struct {
	Profile    *string             `mapstructure:"profile" required:"false" chain:"config.source.profile"`
	RoleArn    *string             `mapstructure:"roleArn" required:"false" chain:"config.source.roleArn"`
	AmiFilters map[*string]*string `mapstructure:"amiFilters" required:"true" chain:"config.source.amiFilters"`
}

type SrcAmiInfo struct {
	CredsInfo       map[string]string
	AccountID       *string
	Filters         []*ec2.Filter
	RegionAmiErrMap RegionAmiErrMap
}

type RawAccountRegionMapping struct {
	AccountID                 int                 `mapstructure:"accountId" required:"true" chain:"config.target.accountRegionMappingList.N.accountId"`
	Profile                   *string             `mapstructure:"profile" required:"false" chain:"config.target.accountRegionMappingList.N.profile"`
	RoleArn                   *string             `mapstructure:"roleArn" required:"false" chain:"config.target.accountRegionMappingList.N.roleArn"`
	Regions                   []*string           `mapstructure:"regions" required:"true" chain:"config.target.accountRegionMappingList.N.regions"`
	AddExtraTags              map[*string]*string `mapstructure:"addExtraTags" required:"false" chain:"config.target.accountRegionMappingList.N.addExtraTags"`
	CopyTagsAcrossAccounts    bool                `mapstructure:"copyTagsAcrossAccounts" required:"false"`
	AddCreateVolumePermission bool                `mapstructure:"addCreateVolumePermission" required:"false"`
}

type AccountRegionMapping struct {
	CopyTags  bool
	AddCVP    bool
	Tags      []*ec2.Tag
	CredsInfo map[string]string
	AccountID *string
	Image     *ec2.Image
	Regions   []*string
}

type Target struct {
	AccountRegionMappingList    []RawAccountRegionMapping `mapstructure:"accountRegionMappingList" required:"true" chain:"config.target.accountRegionMappingList"`
	CopyTagsAcrossAccounts      bool                      `mapstructure:"copyTagsAcrossAccounts" required:"false"`
	CommonRegions               []*string                 `mapstructure:"commonRegions" required:"false" chain:"config.target.commonRegions"`
	AddCreateVolumePermission   bool                      `mapstructure:"addCreateVolumePermission" required:"false"`
	ModAccountRegionMappingList []AccountRegionMapping    `mapstructure:"-"`
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

type (
	AmiInfo struct {
		Ami   *ec2.Image
		Error error
	}
	RegionAmiErrMap map[*string]AmiInfo
	Flag            struct {
		AddCVP bool
	}
	AccountFlagMap map[*string]Flag
)

// Run applies the resource specific configuration.
func (r *Resource) Run() error {
	if err := r.Config.apply(); err != nil {
		return err
	}

	return nil
}

func (t Target) getTargetRegions() []*string {
	regions := make([]*string, 0)

	for _, rawAccountRegionMapping := range t.AccountRegionMappingList {
		regions = append(regions, rawAccountRegionMapping.Regions...)
	}

	regions = append(regions, t.CommonRegions...)

	return regions
}
