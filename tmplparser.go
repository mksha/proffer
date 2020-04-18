package main

import (
	"log"
	"os"
	"text/template"
)

var (
	fm = template.FuncMap{
		"env":     getEnv,
		"default": setDefaultValue,
	}
	tmplPtr *template.Template
)

type (
	Env  map[string]string
	Data struct {
		Env Env
	}
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

// func getOSEnvVars() map[string]string {
// 	envMap := make(map[string]string)
// 	for _, kvStr := range os.Environ() {
// 		kvSlice := strings.Split(kvStr, "=")
// 		envMap[kvSlice[0]] = kvSlice[1]
// 	}
// 	return envMap
// }

func parseTmpl(dsc string){
	log.Println(" ********************** Start: Template Parsing ************************************")
	fileInfo, err := os.Stat(dsc)
	if os.IsNotExist(err) {
		log.Fatalln(err)
	}
	log.Println(" DSC Found At :-> ", dsc)

	dscName := fileInfo.Name()
	// Create template object for given dsc
	tmplPtr = template.Must(template.New(dscName).Funcs(fm).ParseFiles(dsc))
	// Create output file for storing parsed template content
	fPtr, err := os.Create("output.yml")
	if err != nil {
		log.Fatalln(err)
	}
	// Parse given dsc template
	err = tmplPtr.Execute(fPtr, nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(" ************************ End: Template Parsing ************************************")
}

