package main

import (
	"fmt"
	"log"

	"example.com/amidist/command"
	"example.com/amidist/parser"
)

func execute(dsc string) {
	executeResources(dsc)
}

func executeResources(dsc string) {
	c := prepareDataStore(dsc)
	// fmt.Println(c.Data)
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

func prepareDataStore(dsc string) parser.Config {
	log.Println(" ****************** Start: Template Parsing *********************")

	parsedTemplateFileName := parser.ParseTemplate(dsc)
	config := parser.UnmarshalYaml(parsedTemplateFileName)

	log.Println(" ****************** End: Template Parsing *********************")

	return config
}
