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
	"bufio"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/utils"
)

var cfgFile string
var lwCliInst instance.Client
var useContext string

var RootCmd = &cobra.Command{
	Use:   "lw",
	Short: "CLI interface for LiquidWeb",
	Long: `CLI interface for LiquidWeb.

Command line interface for interacting with LiquidWeb services via
LiquidWeb's Public API.

If this is your first time running, you will need to setup at least
one auth context. An auth context contains authentication data for
accessing your LiquidWeb account. As such one auth context represents
one LiquidWeb account. You can have multiple auth contexts defined.

To setup your first auth context, you run 'auth init'. For further
information on auth contexts, be sure to checkout 'help auth' for a
list of capabilities.

As always, consult the various subcommands for specific features and
capabilities.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.liquidweb-cli.yaml)")
	RootCmd.PersistentFlags().StringVar(&useContext, "use-context", "", "forces current context, without persisting the context change")
}

func initConfig() {
	vp := viper.New()
	if cfgFile != "" {
		// Use config file from the flag.
		vp.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			lwCliInst.Die(err)
		}

		// Search config in home directory with name ".liquidweb-cli" (without extension).
		vp.AddConfigPath(home)
		vp.SetConfigName(".liquidweb-cli")
	}

	vp.AutomaticEnv()
	if err := vp.ReadInConfig(); err != nil {
		utils.PrintYellow("no config\n")
	}

	if useContext != "" {
		if err := instance.ValidateContext(useContext, vp); err != nil {
			utils.PrintRed("error using auth context:\n\n")
			fmt.Printf("%s\n\n", err)
			os.Exit(1)
		}
		vp.Set("liquidweb.api.current_context", useContext)
	}

	var lwCliInstErr error
	lwCliInst, lwCliInstErr = instance.New(vp)
	if lwCliInstErr != nil {
		lwCliInst.Die(lwCliInstErr)
	}
}

func dialogDesctructiveConfirmProceed() (proceed bool) {

	var haveConfirmationAnswer bool
	utils.PrintTeal("Tip: Avoid future confirmations by passing --force\n\n")

	for !haveConfirmationAnswer {
		utils.PrintRed("This is a destructive operation. Continue (yes/[no])?: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		answer := scanner.Text()

		if answer != "" && answer != "yes" && answer != "no" {
			utils.PrintYellow("invalid input.\n")
			continue
		}

		haveConfirmationAnswer = true
		if answer == "no" || answer == "" {
			proceed = false
		} else if answer == "yes" {
			proceed = true
		}
	}

	return
}
