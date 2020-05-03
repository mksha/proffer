package shareami

import (
	"log"
	"os"
	"fmt"

	clog "example.com/proffer/common/clogger"
	"github.com/mitchellh/mapstructure"
	awscommon "example.com/proffer/resources/aws/common"
)

var (
	clogger = clog.New(os.Stdout, "aws-shareami | ", log.Lmsgprefix)
)


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
	Source awscommon.RawSrcAmiInfo `yaml:"source"`
	Target Target `yaml:"target"`
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

	return nil
}

func (r *Resource) Run() error {
	source := r.Config.Source
	target := r.Config.Target
	
	fmt.Println(source)
	fmt.Println(target)


	if err := shareAmi(); err != nil {
		return err
	}

	return nil
}
