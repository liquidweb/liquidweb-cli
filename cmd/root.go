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
	"os"
	"regexp"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/config"
	"github.com/liquidweb/liquidweb-cli/flags/defaults"
	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/utils"
)

var cfgFile string
var lwCliInst *instance.Client
var useContext string

var rootCmd = &cobra.Command{
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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.liquidweb-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&useContext, "use-context", "", "forces current context, without persisting the context change")
}

func setConfigArgs() {
	config.UseContextArg = useContext
	config.ConfigFileArg = cfgFile

	if config.UseContextArg == "" {
		osArgsForContext(regexp.MustCompile(`use-context\s+[A-z]+`))
		if config.UseContextArg == "" {
			osArgsForContext(regexp.MustCompile(`use-context=[A-z]+`))
		}
	}
}

// this gets called early on before cobra fully initializes, so have to
// parse os.Args directly.
func osArgsForContext(re *regexp.Regexp) {
	var searchStr string
	for _, str := range os.Args {
		searchStr = searchStr + " " + str
	}

	delimiter := " "
	if strings.Contains(searchStr, "=") {
		delimiter = "="
	}

	slice := strings.Split(re.FindString(searchStr), delimiter)
	if len(slice) > 1 && slice[1] != "" {
		config.UseContextArg = slice[1]
		config.CurrentContext = slice[1]
	}
}

func defaultFlag(flag string, defaultValueList ...interface{}) (value interface{}) {
	// calling config.InitConfig() here so default context gets set
	_, _ = config.InitConfig()
	setConfigArgs()
	value = defaults.GetOrNag(flag)
	if len(defaultValueList) > 0 && value == nil {
		value = defaultValueList[0]
	}
	return
}

func initConfig() {
	vp, err := config.InitConfig()
	if err != nil {
		lwCliInst.Die(err)
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

	f := func(d prompt.Document) []prompt.Suggest {
		s := []prompt.Suggest{
			{Text: "yes", Description: "I understand continue"},
			{Text: "no", Description: "I would like to cancel"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}

	for !haveConfirmationAnswer {
		utils.PrintRed("This is a destructive operation. Continue? ")
		answer := strings.ToLower(prompt.Input("> ", f, prompt.OptionShowCompletionAtStart()))
		if answer == "yes" || answer == "no" || answer == "" {
			haveConfirmationAnswer = true
			if answer == "yes" {
				proceed = true
			}
		}
	}

	return
}
