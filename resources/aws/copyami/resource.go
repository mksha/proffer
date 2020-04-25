package copyami

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mitchellh/mapstructure"
)

// type AmiFilter struct {
// 	ID   *string `yaml:"id"`
// 	Name *string `yaml:"name"`
// }

type Source struct {
	Environment string              `yaml:"environment"`
	Region      *string             `yaml:"region"`
	AmiFilters  map[*string]*string `yaml:"amiFilters"`
}

type Target struct {
	Regions      []string          `yaml:"regions"`
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
	log.Println(" ************************ Preparing Resource **************************** ")

	var c Config

	if err := mapstructure.Decode(rawConfig, &c); err != nil {
		return err
	}

	r.Config = c

	return nil
}

func (r *Resource) Run() error {
	fmt.Println(r.Config)

	source := r.Config.Source
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
		Regions: []*string{
			aws.String("ap-southeast-1"),
		},
	}
	copyAmi(srcAmiInfo, targetInfo)
	return nil
}
