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
	"log"
	"path/filepath"

	"example.com/proffer/command"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply proffer configuration",
	Long: `Apply command is used to apply the proffer configuration and distribute the cloud image
in between multiple regions and with multiple accounts.`,
	Run: applyConfig,
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func applyConfig(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatalln("Proffer Configuration file missing: Pls pass proffer config file to apply")
	}

	cfgFileAbsPath, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatal(err)
	}

	executeResources(cfgFileAbsPath)
}

func executeResources(dsc string) {

	c, err := parseConfig(dsc)
	if err != nil {
		fmt.Println("Unable to parse configuration file")
	}

	resources := command.Resources
	for _, rawResource := range c.RawResources {
		resource, ok := resources[rawResource.Type]
		if !ok {
			log.Fatalf("InvalidResource: Resource %s Not Found", rawResource.Type)
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
