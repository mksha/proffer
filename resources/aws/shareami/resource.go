package shareami

import (
	"log"
	"os"
	"strconv"

	clog "example.com/proffer/common/clogger"
	awscommon "example.com/proffer/resources/aws/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mitchellh/mapstructure"
)

var (
	clogger = clog.New(os.Stdout, "aws-shareami | ", log.Lmsgprefix)
)

type RawAccountRegionMapping struct {
	AccountID              int                 `yaml:"accountId"`
	Profile                *string             `yaml:"profile"`
	RoleArn                *string             `yaml:"roleArn"`
	Regions                []*string           `yaml:"regions"`
	AddExtraTags           map[*string]*string `yaml:"addExtraTags"`
	CopyTagsAcrossAccounts bool                `yaml:"copyTagsAcrossAccounts"`
}

type AccountRegionMapping struct {
	CopyTags  bool
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
	ModAccountRegionMappingList []AccountRegionMapping    `yaml:"-"`
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

	r.prepareAccountRegionMappingList()
	r.Config.Target.setCommonPropertiesIfAny()

	return nil
}

func (r *Resource) Run() error {

	if err := r.Config.apply(); err != nil {
		return err
	}

	return nil
}

func (r *Resource) prepareAccountRegionMappingList() {
	accountRegionMappingList := make([]AccountRegionMapping, 0)

	for _, rawAccountRegionMapping := range r.Config.Target.AccountRegionMappingList {
		accountRegionMapping := AccountRegionMapping{
			CopyTags:  rawAccountRegionMapping.CopyTagsAcrossAccounts,
			Tags:      awscommon.FormEc2Tags(rawAccountRegionMapping.AddExtraTags),
			Regions:   rawAccountRegionMapping.Regions,
			AccountID: aws.String(strconv.Itoa(rawAccountRegionMapping.AccountID)),
			CredsInfo: make(map[string]string, 2),
		}

		if rawAccountRegionMapping.RoleArn != nil {
			accountRegionMapping.CredsInfo["getCredsUsing"] = "roleArn"
			accountRegionMapping.CredsInfo["roleArn"] = *rawAccountRegionMapping.RoleArn
		} else if rawAccountRegionMapping.Profile != nil {
			accountRegionMapping.CredsInfo["getCredsUsing"] = "profile"
			accountRegionMapping.CredsInfo["profile"] = *rawAccountRegionMapping.Profile
		}

		accountRegionMappingList = append(accountRegionMappingList, accountRegionMapping)
	}

	r.Config.Target.ModAccountRegionMappingList = accountRegionMappingList
}

func (t *Target) getTargetAccounts() []*string {
	accounts := make([]*string, 0)

	for _, rawAccountRegionMapping := range t.AccountRegionMappingList {
		accounts = append(accounts, aws.String(string(rawAccountRegionMapping.AccountID)))
	}

	return accounts
}

func (t *Target) setCommonPropertiesIfAny() {
	if t.CommonRegions != nil && t.CopyTagsAcrossAccounts {
		for i := 0; i < len(t.ModAccountRegionMappingList); i++ {
			t.ModAccountRegionMappingList[i].Regions = append(t.ModAccountRegionMappingList[i].Regions, t.CommonRegions...)
			t.ModAccountRegionMappingList[i].CopyTags = t.CopyTagsAcrossAccounts
		}
		return
	}

	if t.CommonRegions != nil {
		for i := 0; i < len(t.ModAccountRegionMappingList); i++ {
			t.ModAccountRegionMappingList[i].Regions = append(t.ModAccountRegionMappingList[i].Regions, t.CommonRegions...)
		}
		return
	}

	if t.CopyTagsAcrossAccounts {
		for i := 0; i < len(t.ModAccountRegionMappingList); i++ {
			t.ModAccountRegionMappingList[i].CopyTags = t.CopyTagsAcrossAccounts
		}
		return
	}
}
