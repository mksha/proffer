/*
Copyright © 2020 flashtaken <flashtaken1@gmail.com>

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
	"log"
	"path/filepath"

	"github.com/proffer/command"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command used to apply the given configuration.
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply proffer configuration",
	Long: `Apply command is used to apply the proffer configuration and distribute the cloud image
in between multiple regions and with multiple accounts.`,
	Run: applyConfig,
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

// applyConfig applies the given template configuration.
func applyConfig(cmd *cobra.Command, args []string) {
	// validate template before applying
	clogger.SetPrefix("start-validation| ")
	fmt.Println()
	clogger.Info("Validating template before applying...")
	validateConfig(cmd, args)
	fmt.Println()

	clogger.SetPrefix("start-apply | ")
	clogger.Info("Applying template config...")

	if len(args) == 0 {
		log.Fatalln("Proffer Configuration file missing: Pls pass proffer config file to apply")
	}

	cfgFileAbsPath, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatal(err)
	}

	// apply template
	executeResources(cfgFileAbsPath)
}

// executeResources applies the given resources in given configuration.
func executeResources(dsc string) {
	c, err := parseConfig(dsc)
	if err != nil {
		fmt.Println("Unable to parse configuration file")
	}

	resources := command.Resources

	// apply resources defined in template one by one
	for _, rawResource := range c.RawResources {
		resource, ok := resources[rawResource.Type]
		if !ok {
			clogger.Fatalf("InvalidResourceType: Resource Type '%s' Not Found", rawResource.Type)
		}

		clogger.SetPrefix(rawResource.Type + " | ")
		clogger.Successf("Resource : %s  Status: Started", rawResource.Name)
		clogger.Info("")

		if err := resource.Prepare(rawResource.Config); err != nil {
			clogger.Error(err)
			clogger.Fatalf("Resource : %s  Status: Failed", rawResource.Name)
		}

		if err := resource.Run(); err != nil {
			clogger.Error(err)
			clogger.Fatalf("Resource : %s  Status: Failed", rawResource.Name)
		}

		clogger.Info("")
		clogger.Successf("Resource : %s  Status: Succeeded", rawResource.Name)
		fmt.Println("")
	}
}
