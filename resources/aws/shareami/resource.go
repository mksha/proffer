package shareami

import (
	"log"
	"os"

	clog "example.com/proffer/common/clogger"
	"github.com/mitchellh/mapstructure"
)

var (
	// logger  = log.New(os.Stdout, common.Success("aws-shareami | "), log.Lmsgprefix)
	clogger = clog.New(os.Stdout, "aws-shareami | ", log.Lmsgprefix)
)

type AmiFilters struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type EnvironmentInfo struct {
	Regions   []string          `yaml:"regions"`
	ExtraTags map[string]string `yaml:"extraTags"`
}

type Source struct {
	Environment string     `yaml:"environment"`
	Region      string     `yaml:"region"`
	AmiFilters  AmiFilters `yaml:"amiFilters"`
}

type Target struct {
	EnvironmentRegionMapping map[string]EnvironmentInfo `yaml:"environmentRegionMapping"`
	CommonRegions            []string                   `yaml:"commonRegions"`
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
