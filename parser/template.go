package parser

import (
	"fmt"
	"os"
	"reflect"
	"text/template"
)

func getEnv(name string) string {
	return os.Getenv(name)
}

func getVar(name string) string {
	value, ok := dynamicVars[name]
	if !ok {
		clogger.Debugf("Variable '%s' not found in variable file. Will check default vars in template file.", name)

		defaultValue, ok := defaultVars["vars"][name]
		if !ok {
			clogger.Fatalf("Variable '%s' not found in default vars [vars property of proffer template file].", name)
		}

		value = defaultValue
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Slice:
		s := "["
		for _, v := range value.([]interface{}) {
			if s == "[" {
				s += v.(string)
			} else {
				s += ", " + v.(string)
			}
		}

		s += "]"

		return s
	case reflect.Map:
		s := "{"

		for k, v := range value.(map[interface{}]interface{}) {
			if pair := k.(string) + ": " + v.(string); s == "{" {
				s += pair
			} else {
				s += ", " + pair
			}
		}

		s += "}"

		return s
	default:
		return fmt.Sprintf("%v", value)
	}
}

func setDefaultValue(givenValue, currentValue string) string {
	if currentValue == "" {
		return givenValue
	}

	return currentValue
}

// ParseTemplate can be used to parse the given go template.
func ParseTemplate(dsc string) (string, error) {
	var (
		fm = template.FuncMap{
			"env":     getEnv,
			"default": setDefaultValue,
			"var":     getVar,
		}
		// Create output file for storing parsed template content
		// parsedTemplateFileName = uuid.New().String() + ".yml"
		parsedTemplateFileName = "output.yml"
	)

	fileInfo, err := os.Stat(dsc)

	if os.IsNotExist(err) {
		return "", err
	}

	clogger.Debug("")
	clogger.Debug("Template Found At :-> ", dsc)

	dscName := fileInfo.Name()
	// Create template object for given dsc
	tmplPtr, err := template.New(dscName).Funcs(fm).ParseFiles(dsc)
	if err != nil {
		return "", err
	}

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
