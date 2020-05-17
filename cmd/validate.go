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
		Run:   validateConfig,
	}
	// clogger = clog.New(os.Stdout, "config-validation | ", log.Lmsgprefix)
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

func validateConfig(cmd *cobra.Command, args []string) {
	config := getTempConfigOnValidSyntax(args)
	validateResources(config)
}

func getTempConfigOnValidSyntax(args []string) parser.TemplateConfig {
	var config parser.TemplateConfig

	clogger.SetPrefix("validate-syntax | ")

	if len(args) == 0 {
		clogger.Fatal("Proffer template file is missing: Pls pass proffer template file to apply")
	}

	cfgFileAbsPath, err := filepath.Abs(args[0])
	if err != nil {
		clogger.Fatal(err)
	}

	if config, err = parseConfig(cfgFileAbsPath); err != nil {
		clogger.Errorf("InvalidTemplate: Unable to parse proffer template file: '%s'", cfgFileAbsPath)
		clogger.Fatal(err)
	}

	clogger.Success("Template syntax is valid.")

	return config
}

func validateResources(c parser.TemplateConfig) {
	resources := command.Resources

	clogger.SetPrefix("validate-config | ")

	// check if the resource list is empty
	if len(c.RawResources) == 0 {
		clogger.Fatal("NoResourceFound: 'resources' list is empty")
	}

	for index, rawResource := range c.RawResources {
		if validator.IsZero(rawResource) {
			clogger.Fatalf("Empty resource found in list 'resources' at index: [%v]", index+1)
		}

		if errs := validator.CheckRequiredFieldsInStruct(rawResource); len(errs) != 0 {
			// check if the resource name is empty
			if validator.IsZero(rawResource.Name) {
				clogger.Errorf("Missing/Empty key(s) in the resource number: [%v]", index+1)
			} else {
				clogger.Errorf("Missing/Empty key(s) in the resource: [%v]", rawResource.Name)
			}

			clogger.Fatal(errs)
		}

		// check if the given resource is valid resource type
		resource, ok := resources[rawResource.Type]
		if !ok {
			clogger.Fatalf("Invalid resource type [%s] found in Resource: [%s]", rawResource.Type, rawResource.Name)
		}

		if err := resource.Validate(rawResource); err != nil {
			clogger.Fatal(err)
		}
	}

	clogger.Success("Template config is valid.")
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
