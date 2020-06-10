package shareami

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/ec2"
	clog "github.com/proffer/common/clogger"
	awscommon "github.com/proffer/resources/aws/common"
	"gopkg.in/yaml.v2"
)

var (
	clogger = clog.New(os.Stdout, "aws-shareami | ", log.Lmsgprefix)
)

// RawSrcAmiInfo represents raw source ami information for aws-shareami resource.
type RawSrcAmiInfo struct {
	Profile    *string             `mapstructure:"profile" required:"false" chain:"config.source.profile"`
	RoleArn    *string             `mapstructure:"roleArn" required:"false" chain:"config.source.roleArn"`
	AmiFilters map[*string]*string `mapstructure:"amiFilters" required:"true" chain:"config.source.amiFilters"`
}

// SrcAmiInfo represents source ami information for aws-shareami resource.
type SrcAmiInfo struct {
	CredsInfo       map[string]string
	AccountID       *string
	AccountAlias    *string
	Filters         []*ec2.Filter
	RegionAmiErrMap RegionAmiErrMap
	RegionalRecord  map[*string]awscommon.AmiMeta
	AccountRecord   map[*string]AccountImage
}

// RawAccountRegionMapping represents raw mapping of accounts and regions.
type RawAccountRegionMapping struct {
	AccountID                 int                 `mapstructure:"accountId" required:"true" chain:"config.target.accountRegionMappingList.N.accountId"`
	AccountAlias              *string             `mapstructure:"accountAlias" required:"false" chain:"config.target.accountRegionMappingList.N.accountAlias"`
	Profile                   *string             `mapstructure:"profile" required:"false" chain:"config.target.accountRegionMappingList.N.profile"`
	RoleArn                   *string             `mapstructure:"roleArn" required:"false" chain:"config.target.accountRegionMappingList.N.roleArn"`
	Regions                   []*string           `mapstructure:"regions" required:"true" chain:"config.target.accountRegionMappingList.N.regions"`
	AddExtraTags              map[*string]*string `mapstructure:"addExtraTags" required:"false" chain:"config.target.accountRegionMappingList.N.addExtraTags"`
	CopyTagsAcrossAccounts    bool                `mapstructure:"copyTagsAcrossAccounts" required:"false"`
	AddCreateVolumePermission bool                `mapstructure:"addCreateVolumePermission" required:"false"`
}

// AccountRegionMapping represents mapping of accounts and regions.
type AccountRegionMapping struct {
	CopyTags     bool
	AddCVP       bool
	Tags         []*ec2.Tag
	CredsInfo    map[string]string
	AccountID    *string
	AccountAlias *string
	Image        *ec2.Image
	Regions      []*string
}

// Target represents target configuration for aws-shareami resource.
type Target struct {
	AccountRegionMappingList    []RawAccountRegionMapping `mapstructure:"accountRegionMappingList" required:"true" chain:"config.target.accountRegionMappingList"`
	CopyTagsAcrossAccounts      bool                      `mapstructure:"copyTagsAcrossAccounts" required:"false"`
	CommonRegions               []*string                 `mapstructure:"commonRegions" required:"false" chain:"config.target.commonRegions"`
	AddCreateVolumePermission   bool                      `mapstructure:"addCreateVolumePermission" required:"false"`
	ModAccountRegionMappingList []AccountRegionMapping    `mapstructure:"-"`
}

// Config represents configuration needed for aws-shareami resource.
type Config struct {
	Source     RawSrcAmiInfo `mapstructure:"source" required:"true" chain:"config.source"`
	Target     Target        `mapstructure:"target" required:"true" chain:"config.target"`
	SrcAmiInfo SrcAmiInfo    `mapstructure:"-"`
}

// SrcImage represent the source ami information used for inventory generation.
type SrcImage struct {
	AmiFilters map[*string]*string   `yaml:"identifiers"`
	Account    awscommon.AccountMeta `yaml:"account"`
}

// AccountImage represents the data of images at account level.
type AccountImage struct {
	AccountAlias *string                       `yaml:"accountAlias"`
	Regions      map[*string]awscommon.AmiMeta `yaml:"regions"`
}

// Record represents the inventory record for aws-shareami resource.
type Record struct {
	SrcImage             SrcImage                 `yaml:"sourceImage"`
	TargetAccountsImages map[*string]AccountImage `yaml:"targetAccountsImages"`
}

// Resource represents aws-shareami resource type.
type Resource struct {
	Name   *string `required:"true"`
	Type   *string `required:"true"`
	Config Config  `mapstructure:"config" required:"true"`
	Record Record  `mapstructure:"-"`
}

type (
	// AmiInfo represent metadata of Amazon Machine Image.
	AmiInfo struct {
		Ami   *ec2.Image
		Error error
	}
	// RegionAmiErrMap represents the mapping of region, ami and error.
	RegionAmiErrMap map[*string]AmiInfo
	// Flag represents the different flags that can be passed to accoutRegion mapping config.
	Flag struct {
		AddCVP bool
	}
	// AccountFlagMap represents the flag mapping.
	AccountFlagMap map[*string]Flag
)

// Run applies the resource specific configuration.
func (r *Resource) Run() error {
	if err := r.Config.apply(); err != nil {
		return err
	}

	r.Record.SrcImage.AmiFilters = r.Config.Source.AmiFilters
	r.Record.SrcImage.Account = awscommon.AccountMeta{
		ID:    r.Config.SrcAmiInfo.AccountID,
		Alias: r.Config.SrcAmiInfo.AccountAlias,
	}

	r.Record.TargetAccountsImages = make(map[*string]AccountImage)
	for account, accountImage := range r.Config.SrcAmiInfo.AccountRecord {
		r.Record.TargetAccountsImages[account] = accountImage
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

// GenerateInventory generates the distribution inventory for aws-shareami resource.
func (r *Resource) GenerateInventory() ([]byte, error) {
	inventoryRecord := struct {
		ResourceName *string `yaml:"resourceName"`
		Output       Record  `yaml:"output"`
	}{
		ResourceName: r.Name,
		Output:       r.Record,
	}

	bs, err := yaml.Marshal(inventoryRecord)

	if err != nil {
		return nil, err
	}

	return bs, nil
}
