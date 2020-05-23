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
   "regexp"
   "testing"

	"github.com/spf13/cobra"
)

func Test_applyCmd(t *testing.T) {
	root := rootCmd

	root.ResetCommands()

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply proffer configuration",
		Long: `Apply command is used to apply the proffer configuration and distribute the cloud image
in between multiple regions and with multiple accounts.`,
		RunE: applyConfig,
	}
	root.AddCommand(applyCmd)

	for n := range applyCmdTestCases {
		tt := applyCmdTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeCommand(root, tt.args...)
			matched, _ := regexp.MatchString(tt.want, output)
			if !matched {
				t.Errorf("Unexpected output: %v, want: %v", output, tt.want)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
