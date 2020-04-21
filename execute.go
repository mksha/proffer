package main

import (
	"fmt"
	"log"

	"example.com/amidist/command"
)

func execute(dsc string) {
	executeResources(dsc)
}

func executeResources(dsc string) {
	c := prepareDataStore(dsc)
	resources := command.Resources
	for _, rawResource := range c.RawResources {
		fmt.Println(rawResource.Prepare(rawResource.Config))
		resource, ok := resources[rawResource.Type]
		if !ok {
			log.Fatalf(" InvalidResource: Resource %s Not Found", rawResource.Type)
		}
		// c.Resources[rawResource.Type] = resource

		fmt.Println(resource)
		fmt.Printf(" ******** before ***********\n %T\n", resource)
		fmt.Println(resource.Prepare(rawResource.Config))
		fmt.Printf(" ******** after ************\n %T\n", resource)
		// resource.Run()
	}
}

func prepareDataStore(dsc string) config {
	log.Println(" ****************** Start: Template Parsing *********************")

	parsedTemplateFileName := parseTemplate(dsc)
	config := unmarshalYaml(parsedTemplateFileName)

	log.Println(" ****************** End: Template Parsing *********************")

	return config
}
