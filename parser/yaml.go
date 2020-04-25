package parser

import (
	"log"
	"os"

	"example.com/amidist/components"
	"gopkg.in/yaml.v3"
)

type Env struct {
	AccountID     int    `yaml:"accountId"`
	GetCredsUsing string `yaml:"getCredsUsing"`
	IamRoleName   string `yaml:"iamRoleName"`
}

type Environments map[string]*Env

type Data struct {
	Environments Environments `yaml:"environments"`
}

type Resource struct {
	Name   string                 `yaml:"name"`
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config"`
}

type Config struct {
	Data         Data                     `yaml:"data"`
	RawResources []Resource               `yaml:"resources"`
	Resources    components.MapOfResource `yaml:"-"`
}

func UnmarshalYaml(filePath string) Config {
	var c Config

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Fatalln(err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	log.Println(" YAML file found :-> ", filePath)

	data := make([]byte, fileInfo.Size())

	if _, err := file.Read(data); err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}

	return c
}
