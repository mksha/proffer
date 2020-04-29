package copyami

import (
	"log"
	"os"

	clog "example.com/proffer/common/clogger"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mitchellh/mapstructure"
)

const (
	LogPrefix = "aws-copyami | "
)

var (
	// infoLog  = log.New(os.Stdout, common.Success("aws-copyami | "), log.Lmsgprefix)
	// errorLog = log.New(os.Stdout, common.Error("aws-copyami | "), log.Llongfile)
	// errorMsg = log.New(os.Stdout, common.Error("aws-copyami | "), log.Lmsgprefix)
	clogger = clog.New(os.Stdout, "aws-copyami | ", log.Lmsgprefix)
)

type Source struct {
	Environment string              `yaml:"environment"`
	Region      *string             `yaml:"region"`
	AmiFilters  map[*string]*string `yaml:"amiFilters"`
}

type Target struct {
	Regions      []*string         `yaml:"regions"`
	AddExtraTags map[string]string `yaml:"addExtraTags"`
}

type Config struct {
	Source Source `yaml:"source"`
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
	var amiFilters []*ec2.Filter

	for filterName, filterValue := range source.AmiFilters {
		f := &ec2.Filter{
			Name:   filterName,
			Values: []*string{filterValue},
		}
		amiFilters = append(amiFilters, f)
	}

	srcAmiInfo := SrcAmiInfo{
		Region:  source.Region,
		Filters: amiFilters,
	}

	targetInfo := TargetInfo{
		Regions: target.Regions,
	}

	copyAmi(srcAmiInfo, targetInfo)

	return nil
}
