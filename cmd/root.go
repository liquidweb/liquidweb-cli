/*
Copyright © LiquidWeb

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
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/liquidweb/liquidweb-cli/instance"
)

var cfgFile string
var lwCliInst instance.Client

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "liquidweb-cli",
	Short: "CLI interface for LiquidWeb",
	Long: `CLI interface for LiquidWeb.

Command line interface for interacting with LiquidWeb services via
LiquidWebs Public API.

Your public API credentials must be in the liquidweb-cli config
file. By default, this will be $HOME/.liquidweb-cli.yaml. This can
be overrode with the config flag.

Consult the various subcommands for specific features and
capabilities.`,
}

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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.liquidweb-cli.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	vp := viper.New()
	if cfgFile != "" {
		// Use config file from the flag.
		vp.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".liquidweb-cli" (without extension).
		vp.AddConfigPath(home)
		vp.SetConfigName(".liquidweb-cli")
	}

	vp.AutomaticEnv()
	vp.ReadInConfig()

	var lwCliInstErr error
	lwCliInst, lwCliInstErr = instance.New(vp)
	if lwCliInstErr != nil {
		fmt.Printf("Fatal: [%s]\n", lwCliInstErr)
		os.Exit(1)
	}
}
