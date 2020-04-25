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
	"log"
	"path/filepath"

	"example.com/proffer/parser"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate proffer configuration file.",
	Long:  `Validate command is used to validate the proffer configuration file.`,
	Run:   validateConfig,
}

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
	if len(args) == 0 {
		log.Fatalln(" Proffer Configuration file missing: Pls pass proffer config file to apply")
	}

	cfgFileAbsPath, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatal(err)
	}

	if _, err := parseConfig(cfgFileAbsPath); err != nil {
		log.Fatalf(" InvalidConfig: Unable to parse configuration file %s \n", cfgFileAbsPath)
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
