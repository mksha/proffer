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
	// "log"
	"path/filepath"

	"example.com/proffer/command"
	"example.com/proffer/parser"
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func validateConfig(cmd *cobra.Command, args []string) {
	config := getTempConfigOnValidSyntax(args)
	validateResources(config)
}

func getTempConfigOnValidSyntax(args []string) parser.Config {
	var config parser.Config

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

	clogger.Success("Template syntax is valid")

	return config
}

func validateResources(c parser.Config) {
	resources := command.Resources
	clogger.SetPrefix("validate-resources | ")

	// check if the resource list is empty
	if len(c.RawResources) == 0 {
		clogger.Fatal("NoResourceFound: Resource list is empty")
	}

	for index, rawResource := range c.RawResources {

		// check if the resource name is empty
		if rawResource.Name == "" {
			clogger.Fatalf("Property 'name' missing for resource number: %v", index+1)
		}

		// Check if the resource type is empty
		if rawResource.Type == "" {
			clogger.Fatalf("Property 'type' missing for resource: %s", rawResource.Name)
		}

		// check if the given resource is valid resource type
		resource, ok := resources[rawResource.Type]
		if !ok {
			clogger.Fatalf("InvalidResourceType: Resource Type '%s' Not Found", rawResource.Type)
		}

		// check if the given resource config is empty
		if rawResource.Config == nil {
			clogger.Fatalf("Property 'config' missing for resource: %s", rawResource.Name)
		}

		if err := resource.Validate(rawResource.Config); err != nil {
			clogger.Error(err)
		}

		clogger.Success("Template config is valid.")
	}
}

func parseConfig(dsc string) (parser.Config, error) {
	var config parser.Config

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
