package copyami

import (
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"
)

type AmiFilters struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type Source struct {
	Environment string     `yaml:"environment"`
	Region      string     `yaml:"region"`
	AmiFilters  AmiFilters `yaml:"amiFilters"`
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

	return nil
}
