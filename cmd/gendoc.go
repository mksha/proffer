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
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	// Used for flags.
	path string
	// gendocCmd represents the gendoc command.
	gendocCmd = &cobra.Command{
		Use:   "gendoc",
		Short: "Generate proffer markdown documentation",
		Long:  `This command is used to generate the proffer markdown documentation using automation.`,
		Run:   generateDoc,
	}
)

func init() {
	rootCmd.AddCommand(gendocCmd)
	gendocCmd.Flags().StringVarP(&path, "path", "p", "doc", "path for generating proffer documentation. Default is 'doc'")
}

// Generate docs.
func generateDoc(cmd *cobra.Command, args []string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			log.Fatal(err)
		}
	}

	if err := doc.GenMarkdownTree(rootCmd, path); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated proffer documents in markdown format in folder %v.", path)
}
