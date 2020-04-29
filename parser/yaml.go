package parser

import (
	"os"

	"example.com/proffer/components"
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

func UnmarshalYaml(filePath string) (Config, error) {
	var c Config

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return c, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return c, err
	}
	defer file.Close()

	clogger.Debug(" YAML file found :-> ", filePath)

	data := make([]byte, fileInfo.Size())

	if _, err := file.Read(data); err != nil {
		return c, err
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}
