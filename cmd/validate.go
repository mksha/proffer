/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/proffer/command"
	"github.com/proffer/common/validator"
	"github.com/proffer/parser"
	"github.com/spf13/cobra"
)

var (
	// validateCmd represents the validate command
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate proffer configuration file.",
		Long:  `Validate command is used to validate the proffer configuration file.`,
		Run:   runValidateConfig,
	}
	// clogger = clog.New(os.Stdout, "config-validation | ", log.Lmsgprefix)
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidateConfig(cmd *cobra.Command, args []string) {
	clogger.Info()

	if err := validateConfig(cmd, args); err != nil {
		clogger.Fatal(err)
	}
}

func validateConfig(cmd *cobra.Command, args []string) error {
	config, err := getTempConfigOnValidSyntax(args)
	if err != nil {
		return err
	}

	if err := validateResources(config); err != nil {
		return err
	}

	return nil
}

func getTempConfigOnValidSyntax(args []string) (parser.TemplateConfig, error) {
	var config parser.TemplateConfig

	clogger.SetPrefix("validate-syntax | ")

	if len(args) == 0 {
		return config, fmt.Errorf("proffer template file is missing: Pls pass proffer template file to apply")
	}

	cfgFileAbsPath, err := filepath.Abs(args[0])
	if err != nil {
		return config, err
	}

	if config, err = parseConfig(cfgFileAbsPath); err != nil {
		clogger.Errorf("InvalidTemplate: Unable to parse proffer template file: '%s'", cfgFileAbsPath)
		return config, err
	}

	clogger.Success("Template syntax is valid.")

	return config, nil
}

func validateResources(c parser.TemplateConfig) error {
	resources := command.Resources

	clogger.SetPrefix("validate-config | ")

	// check if the resource list is empty
	if len(c.RawResources) == 0 {
		return fmt.Errorf("NoResourceFound: 'resources' list is empty")
	}

	for index, rawResource := range c.RawResources {
		if validator.IsZero(rawResource) {
			return fmt.Errorf("empty resource found in list 'resources' at index: [%v]", index+1)
		}

		cs := validator.CustomStruct{Struct: rawResource}
		if errs := validator.CheckRequiredFieldsInStruct(cs); len(errs) != 0 {
			// check if the resource name is empty
			if validator.IsZero(rawResource.Name) {
				clogger.Errorf("Missing/Empty key(s) in the resource number: [%v]", index+1)
			} else {
				clogger.Errorf("Missing/Empty key(s) in the resource: [%v]", rawResource.Name)
			}

			clogger.Error(errs)

			return fmt.Errorf("%v", errs[0])
		}

		// check if the given resource is valid resource type
		resource, ok := resources[rawResource.Type]
		if !ok {
			return fmt.Errorf("invalid resource type [%s] found in Resource: [%s]", rawResource.Type, rawResource.Name)
		}

		if err := resource.Validate(rawResource); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	clogger.Success("Template config is valid.")

	return nil
}

func parseConfig(dsc string) (parser.TemplateConfig, error) {
	var config parser.TemplateConfig

	parsedTemplateFileName, err := parser.ParseTemplate(dsc)
	if err != nil {
		return config, err
	}

	config, err = parser.UnmarshalYaml(parsedTemplateFileName)
	if err != nil {
		return config, err
	}

	return config, nil
}
