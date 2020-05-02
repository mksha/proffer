package shareami

import (
	"log"
	"os"

	clog "example.com/proffer/common/clogger"
	"github.com/mitchellh/mapstructure"
)

var (
	clogger = clog.New(os.Stdout, "aws-shareami | ", log.Lmsgprefix)
)

type Source struct {
	Profile    string              `yaml:"profile"`
	RoleArn    string              `yaml:"roleArn"`
	Region     string              `yaml:"region"`
	AmiFilters map[*string]*string `yaml:"amiFilters"`
}

type AccountRegionMapping struct {
	AccountID              int                 `yaml:"accountId"`
	Profile                string              `yaml:"profile"`
	RoleArn                string              `yaml:"roleArn"`
	Regions                []*string           `yaml:"regions"`
	AddExtraTags           map[*string]*string `yaml:"addExtraTags"`
	CopyTagsAcrossAccounts bool                `yaml:"copyTagsAcrossAccounts"`
}

type Target struct {
	AccountRegionMappingList []AccountRegionMapping `yaml:"accountRegionMapping"`
	CopyTagsAcrossAccounts   bool                   `yaml:"copyTagsAcrossAccounts"`
	CommonRegions            []*string              `yaml:"commonRegions"`
}

type Config struct {
	Source Source `yaml:"source"`
	Target Target `yaml:"target"`
}

type Resource struct {
	Config Config
}

func (r *Resource) Prepare(rawConfig map[string]interface{}) error {
	var c Config

	if err := mapstructure.Decode(rawConfig, &c); err != nil {
		return err
	}

	r.Config = c

	return nil
}

func (r *Resource) Run() error {
	clogger.Info(r.Config)

	return nil
}
