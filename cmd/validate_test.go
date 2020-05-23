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
	"bytes"
	"reflect"
	"regexp"
	"testing"

	"github.com/spf13/cobra"
)

func Test_parseConfig(t *testing.T) {
	for n := range tests {
		tt := tests[n]
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConfig(tt.dsc)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTempConfigOnValidSyntax(t *testing.T) {
	for n := range tests {
		tt := tests[n]
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTempConfigOnValidSyntax([]string{tt.dsc})
			if (err != nil) != tt.wantErr {
				t.Errorf("getTempConfigOnValidSyntax() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTempConfigOnValidSyntax() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateResources(t *testing.T) {
	for n := range tests {
		tt := tests[n]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := validateResources(tt.want); (err != nil) != tt.wantErr {
				t.Errorf("validateResources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}
func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func Test_validateCmd(t *testing.T) {
	root := rootCmd

	root.ResetCommands()

	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate proffer configuration file.",
		Long:  `Validate command is used to validate the proffer configuration file.`,
		RunE:  validateConfig,
	}
	root.AddCommand(validateCmd)
	root.SilenceUsage = true

	for n := range validaCmdTestCases {
		tt := validaCmdTestCases[n]
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
