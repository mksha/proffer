package common

import (
	"fmt"
	"log"

	"example.com/proffer/command"
	"example.com/proffer/parser"
)

func ExecuteResources(dsc string) {
	c, err := ParseConfig(dsc)
	if err != nil {
		fmt.Println(" Unable to parse configuration file")
	}

	resources := command.Resources
	for _, rawResource := range c.RawResources {
		resource, ok := resources[rawResource.Type]
		if !ok {
			log.Fatalf(" InvalidResource: Resource %s Not Found", rawResource.Type)
		}

		fmt.Printf(" ******************** Executing Resource : %s ************************* \n", rawResource.Name)

		if err := resource.Prepare(rawResource.Config); err != nil {
			log.Fatalln(err)
		}

		if err := resource.Run(); err != nil {
			log.Fatalln(err)
		}
	}
}

func ParseConfig(dsc string) (parser.Config, error) {
	log.Println(" ****************** Start: Template Parsing *********************")
	var config parser.Config

	parsedTemplateFileName, err := parser.ParseTemplate(dsc)
	if err != nil {
		return config, err
	}

	config, err = parser.UnmarshalYaml(parsedTemplateFileName)
	if err != nil {
		return config, err
	}

	log.Println(" ****************** End: Template Parsing *********************")

	return config, nil
}
