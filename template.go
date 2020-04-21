package main

import (
	"log"
	"os"
	"text/template"
)

func getEnv(name string) string {
	return os.Getenv(name)
}

func setDefaultValue(givenValue, currentValue string) string {
	if currentValue == "" {
		return givenValue
	}

	return currentValue
}

func parseTemplate(dsc string) string {
	var (
		fm = template.FuncMap{
			"env":     getEnv,
			"default": setDefaultValue,
		}
	)

	fileInfo, err := os.Stat(dsc)

	if os.IsNotExist(err) {
		log.Fatalln(err)
	}

	log.Println(" DSC Found At :-> ", dsc)

	dscName := fileInfo.Name()
	// Create template object for given dsc
	tmplPtr := template.Must(template.New(dscName).Funcs(fm).ParseFiles(dsc))
	// Create output file for storing parsed template content
	parsedTemplateFileName := "output.yml"

	file, err := os.Create(parsedTemplateFileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// Parse given dsc template
	err = tmplPtr.Execute(file, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return parsedTemplateFileName
}
