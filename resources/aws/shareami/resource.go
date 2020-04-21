package shareami

import (
	"fmt"

	"example.com/amidist/config"
)

type AmiFilters map[string]struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type EnvRegionMap map[string]struct {
	Regions      []string          `yaml:"regions"`
	AddExtraTags map[string]string `yaml:"extraTags"`
}

type Source struct {
	Environment string     `yaml:"environment"`
	Region      string     `yaml:"region"`
	AmiFilters  AmiFilters `yaml:"amiFilters"`
}

type Target struct {
	EnvRegionMapList []EnvRegionMap `yaml:"environmentRegionMapping"`
	CommonRegions    []string       `yaml:"commonRegions"`
}

type Config struct {
	Source Source `yaml:"source"`
	Target Target `yaml:"target"`
}

type Resource struct {
	config Config
}

func (r Resource) Prepare(rawConfig map[string]interface{}) string {
	fmt.Println("I am in side of prepare")
	fmt.Println(r)
	var c Config
	config.Decode(rawConfig, c, "aws-shareami")
	fmt.Println(c)
	return fmt.Sprintf("Hi from Prepare of type %T\n", r)
}

func (r Resource) Run() string {
	return fmt.Sprintf("Hi from Run of type %T\n", r)
}
