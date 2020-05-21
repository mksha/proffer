package parser

import (
	"os"

	"gopkg.in/yaml.v2"
)

func UnmarshalYaml(filePath string) (TemplateConfig, error) {
	var c TemplateConfig

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return c, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return c, err
	}
	defer file.Close()

	data := make([]byte, fileInfo.Size())

	if _, err := file.Read(data); err != nil {
		return c, err
	}

	err = yaml.UnmarshalStrict(data, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}
