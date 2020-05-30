/*
Copyright © 2020 mohit-kumar-sharma <flashtaken1@gmail.com>

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
	"errors"
	"os"
	"path/filepath"

	"github.com/lithammer/dedent"
	"github.com/proffer/command"
	"github.com/proffer/common/validator"
	"github.com/proffer/parser"
	"github.com/spf13/cobra"
)

var (
	validateExamples = dedent.Dedent(`
		$ proffer validate [flags] TEMPLATE
		$ proffer validate proffer.yml
		$ proffer validate -d proffer.yml`)

	// validateCmd represents the validate command
	validateCmd = &cobra.Command{
		Use:     "validate",
		Short:   "Validate proffer configuration file.",
		Long:    `Validate command is used to validate the proffer configuration file syntax and configuration itself.`,
		Example: validateExamples,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("proffer config file is missing in arguments, pls pass config file to validate")
			}
			return nil
		},
		Run: validateConfig,
	}
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

// validates the config passed to the tool.
func validateConfig(cmd *cobra.Command, args []string) {
	config := getTempConfigOnValidSyntax(args)
	validateResources(config)

	// cleanup temp files.
	if !debug {
		_ = os.Remove("output.yml")
	}
}

// returns the parsed template config if there was no errors throughout parsing.
func getTempConfigOnValidSyntax(args []string) parser.TemplateConfig {
	var config parser.TemplateConfig

	clogger.SetPrefix("validate-syntax | ")

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

// validates the resources specified in the template config.
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

		cs := validator.CustomStruct{Struct: rawResource}
		if errs := validator.CheckRequiredFieldsInStruct(cs); len(errs) != 0 {
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

// parses the given template config and returns the parsed template with error if there was any.
func parseConfig(dsc string) (parser.TemplateConfig, error) {
	var config parser.TemplateConfig

	if dynamicVarsFile != "" {
		if err := parser.UnmarshalDynamicVars(dynamicVarsFile); err != nil {
			return config, err
		}
	}

	if err := parser.UnmarshalDefaultVars(dsc); err != nil {
		return config, err
	}

	parsedTemplateFileName, err := parser.ParseTemplate(dsc)
	if err != nil {
		return config, err
	}

	config, err = parser.UnmarshalYaml(parsedTemplateFileName)
	if err != nil {
		return config, err
	}

	clogger.Debugf("Parsed config can be found at: %s", parsedTemplateFileName)

	return config, nil
}
