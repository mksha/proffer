package shareami

import (
	"log"
	"os"

	clog "example.com/proffer/common/clogger"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	clogger = clog.New(os.Stdout, "aws-shareami | ", log.Lmsgprefix)
)

type RawSrcAmiInfo struct {
	Profile    *string             `yaml:"profile"`
	RoleArn    *string             `yaml:"roleArn"`
	AmiFilters map[*string]*string `yaml:"amiFilters"`
}

type SrcAmiInfo struct {
	CredsInfo       map[string]string
	AccountID       *string
	Filters         []*ec2.Filter
	RegionAmiErrMap RegionAmiErrMap
}

type RawAccountRegionMapping struct {
	AccountID                 int                 `yaml:"accountId"`
	Profile                   *string             `yaml:"profile"`
	RoleArn                   *string             `yaml:"roleArn"`
	Regions                   []*string           `yaml:"regions"`
	AddExtraTags              map[*string]*string `yaml:"addExtraTags"`
	CopyTagsAcrossAccounts    bool                `yaml:"copyTagsAcrossAccounts"`
	AddCreateVolumePermission bool                `yaml:"addCreateVolumePermission"`
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
	AccountRegionMappingList    []RawAccountRegionMapping `yaml:"accountRegionMappingList"`
	CopyTagsAcrossAccounts      bool                      `yaml:"copyTagsAcrossAccounts"`
	CommonRegions               []*string                 `yaml:"commonRegions"`
	AddCreateVolumePermission   bool                      `yaml:"addCreateVolumePermission"`
	ModAccountRegionMappingList []AccountRegionMapping    `yaml:"-"`
}

type Config struct {
	Source     RawSrcAmiInfo          `yaml:"source"`
	Target     Target                 `yaml:"target"`
	Other      map[string]interface{} `mapstructure:",remain"`
	SrcAmiInfo SrcAmiInfo             `yaml:"-"`
}

type Resource struct {
	Name   *string
	Type   *string
	Config Config `yaml:"config"`
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

func (t Target) getTargetAccountFlagMap() AccountFlagMap {
	accountFlagMap := make(AccountFlagMap, 0)

	for _, accountRegionMapping := range t.ModAccountRegionMappingList {
		account := accountRegionMapping.AccountID
		accountFlagMap[account] = Flag{AddCVP: accountRegionMapping.AddCVP}
	}

	return accountFlagMap
}
