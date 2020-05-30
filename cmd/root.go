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

	"github.com/lithammer/dedent"
	clog "github.com/proffer/common/clogger"
	"github.com/spf13/cobra"
)

var (
	debug           bool
	clogger         = clog.New(os.Stdout, "", log.Lmsgprefix)
	dynamicVarsFile string
)

var (
	profferLong = dedent.Dedent(`
		Proffer is a command-line tool to distribute cloud images in between multiple regions
		and with multiple environments. This tool only needs a yml configuration file with name proffer that defines
		the image distribution operations and actions. Currently AWS cloud is the only supported cloud provider but
		support for other cloud providers can be added via resource plugin.`)

	profferExamples = dedent.Dedent(`
		$ proffer [command] [flags] TEMPLATE
		$ proffer validate proffer.yml
		$ proffer validate -d proffer.yml
		$ proffer apply proffer.yml`)

	// rootCmd represents the base command when called without any subcommands.
	rootCmd = &cobra.Command{
		Use:     "proffer",
		Short:   "Proffer is a cross platform tool to distribute cloud images between multiple regions and environments using yml configuration file.",
		Long:    profferLong,
		Example: profferExamples,
		Version: "0.1.4",
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//	Run: func(cmd *cobra.Command, args []string) { },
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Set debug flag to get detailed logging")
	rootCmd.PersistentFlags().StringVar(&dynamicVarsFile, "var-file", "", "Variable file to pass variable's value, otherwise will use default values")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	setLogLevel()
}

// sets the defaultlog level for all commands and subcommands.
func setLogLevel() {
	// Default level is info, unless debug flag is present
	clog.SetGlobalLogLevel(clog.INFO)

	if debug {
		clog.SetGlobalLogLevel(clog.DEBUG)
	}
}
