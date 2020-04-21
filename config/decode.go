package config

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	// "example.com/amidist/command"
)


func Decode(rawConfig map[string]interface{}, targetConfig interface{}, resourceType string) error {
	// resource := command.Resources[resourceType]
	fmt.Println("Raw config ", rawConfig)
	fmt.Printf("Target config %T\n",targetConfig)
	var md mapstructure.Metadata

	config := &mapstructure.DecoderConfig{
		Metadata: &md,
		Result:   &targetConfig,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	if err := decoder.Decode(rawConfig); err != nil {
		panic(err)
	}

	fmt.Println("I am in side of prepare")
	fmt.Println(targetConfig)

	return nil
}
