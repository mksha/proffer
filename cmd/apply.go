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
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command.
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply proffer configuration",
	Long: `Apply command is used to apply the proffer configuration and distribute the cloud image
in between multiple regions and with multiple accounts.`,
	Run: runApplyConfig,
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

func runApplyConfig(cmd *cobra.Command, args []string) {
	if err := applyConfig(cmd, args); err != nil {
		clogger.Error(err)
	}
}

func applyConfig(cmd *cobra.Command, args []string) error {
	// validate template before applying
	clogger.SetPrefix("start-validation| ")
	fmt.Println()

	if len(args) == 0 {
		return fmt.Errorf("proffer Configuration file missing: Pls pass proffer config file to apply")
	}

	cfgFileAbsPath, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	clogger.Info("Validating template before applying...")

	if err := validateConfig(cmd, args); err != nil {
		return err
	}

	fmt.Println()

	clogger.SetPrefix("start-apply | ")
	clogger.Info("Applying template config...")

	if err := executeResources(cfgFileAbsPath); err != nil {
		return err
	}

	return nil
}

func executeResources(dsc string) error {
	c, err := parseConfig(dsc)
	if err != nil {
		clogger.Errorf("unable to parse configuration file")
		return err
	}

	resources := command.Resources

	for _, rawResource := range c.RawResources {
		resource, ok := resources[rawResource.Type]
		if !ok {
			return fmt.Errorf("invalidResourceType: Resource Type '%s' Not Found", rawResource.Type)
		}

		clogger.SetPrefix(rawResource.Type + " | ")
		clogger.Successf("Resource : %s  Status: Started", rawResource.Name)
		clogger.Info("")

		if err := resource.Prepare(rawResource.Config); err != nil {
			clogger.Errorf("resource : %s  Status: Failed", rawResource.Name)
			return err
		}

		if err := resource.Run(); err != nil {
			clogger.Errorf("resource : %s  Status: Failed", rawResource.Name)
			return err
		}

		clogger.Info("")
		clogger.Successf("Resource : %s  Status: Succeeded", rawResource.Name)
		fmt.Println("")
	}

	return nil
}
