package parser

import (
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

// UnmarshalYaml is used to unmarshal the given yaml file to go struct.
func UnmarshalYaml(filePath string) (TemplateConfig, error) {
	var c TemplateConfig

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c, err
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c, err
	}

	err = yaml.UnmarshalStrict(data, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

// UnmarshalDynamicVars un-marshalls the given dynamic vars file.
func UnmarshalDynamicVars(dynamicVarsPath string) error {
	// Reset dynamicVars global var.
	clogger.Debug("")
	clogger.Debug("Resetting dynamic vars.")

	dynamicVars = DynamicVars{}

	// Check dynamic vars file.
	clogger.Debug("Checking dynamic vars file")

	if _, err := os.Stat(dynamicVarsPath); os.IsNotExist(err) {
		return err
	}

	// Evaluate dynamic vars.
	clogger.Debug("Evaluate dynamic vars")

	evaluatedDynamicVarsFilePath, err := ParseTemplate(dynamicVarsPath)
	if err != nil {
		return err
	}

	// Populate dynamic vars.
	dynamicVarsData, err := ioutil.ReadFile(evaluatedDynamicVarsFilePath)
	if err != nil {
		return err
	}

	clogger.Debug("Populate dynamic vars")

	err = yaml.UnmarshalStrict(dynamicVarsData, &dynamicVars)
	if err != nil {
		return err
	}

	clogger.Debug("Parsed dynamic vars")

	return nil
}

// UnmarshalDefaultVars un-marshalls the default vars from proffer template file.
func UnmarshalDefaultVars(defaultVarsPath string) error {
	// Reset defaultVars global var.
	clogger.Debug("")
	clogger.Debug("Resetting default vars.")

	defaultVars = DefaultVars{}

	// Check default vars file.
	if _, err := os.Stat(defaultVarsPath); os.IsNotExist(err) {
		return err
	}

	// Populate default vars.
	defaultVarsData, err := ioutil.ReadFile(defaultVarsPath)
	if err != nil {
		return err
	}

	// Regex pattern captures "^resources:" from the content
	// because its reserved word to specify resource list so
	// not allowed any other key with same name at top level scope.
	pattern := regexp.MustCompile(`(?m)^resources:$`)
	loc := pattern.FindIndex(defaultVarsData)

	clogger.Debug("Populate default vars if any.")

	err = yaml.UnmarshalStrict(defaultVarsData[:loc[0]], &defaultVars)
	if err != nil {
		return err
	}

	clogger.Debug("Parsed default vars.")

	return nil
}
