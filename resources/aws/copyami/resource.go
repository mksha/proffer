package copyami

// import "example.com/amidist/components"

import (
	"fmt"

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
	config Config `yaml:"config"`
}

func (r Resource) Prepare(rawConfig map[string]interface{}) string {
	fmt.Println("Raw config ", rawConfig)
	var c Config
	var md mapstructure.Metadata

	config := &mapstructure.DecoderConfig{
		Metadata: &md,
		Result:   &c,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	if err := decoder.Decode(rawConfig); err != nil {
		panic(err)
	}
	fmt.Println("I am in side of prepare")
	fmt.Println(c)

	return fmt.Sprintf("Hi from Prepare of type %T\n", r)
}

func (r Resource) Run() string {
	return fmt.Sprintf("Hi from Run of type %T\n", r)
}
