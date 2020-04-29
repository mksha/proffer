package parser

import (
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

func ParseTemplate(dsc string) (string, error) {
	var (
		fm = template.FuncMap{
			"env":     getEnv,
			"default": setDefaultValue,
		}
		// Create output file for storing parsed template content
		parsedTemplateFileName = "output.yml"
	)

	fileInfo, err := os.Stat(dsc)

	if os.IsNotExist(err) {
		return "", err
	}

	clogger.Debug(" DSC Found At :-> ", dsc)

	dscName := fileInfo.Name()
	// Create template object for given dsc
	tmplPtr := template.Must(template.New(dscName).Funcs(fm).ParseFiles(dsc))

	file, err := os.Create(parsedTemplateFileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Parse given dsc template
	err = tmplPtr.Execute(file, nil)
	if err != nil {
		return parsedTemplateFileName, err
	}

	return parsedTemplateFileName, nil
}
