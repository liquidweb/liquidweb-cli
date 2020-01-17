/*
Copyright © 2019 LiquidWeb

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

	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/types/cmd"
)

var authUseContextCmd = &cobra.Command{
	Use:   "use-context",
	Short: "Change your current context",
	Long: `Change your current context.

You must pass the name of the context you wish to change to as an argument. Example:

  auth use-context dev

See also: get-context, get-contexts.

If you've never setup any contexts, check "auth init".`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			os.Exit(1)
		}

		wantedContext := args[0]

		// verify wantedContext is valid
		isValid := false
		contexts := lwCliInst.Viper.GetStringMap("liquidweb.api.contexts")
		for _, contextInter := range contexts {
			var context cmdTypes.AuthContext
			if err := instance.CastFieldTypes(contextInter, &context); err != nil {
				lwCliInst.Die(err)
			}

			if context.ContextName == wantedContext {
				isValid = true
				break
			}
		}

		if !isValid {
			lwCliInst.Die(fmt.Errorf("given context [%s] is not a valid context", wantedContext))
		}

		// looks valid, set
		lwCliInst.Viper.Set("liquidweb.api.current_context", wantedContext)
		if err := lwCliInst.Viper.WriteConfig(); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("context changed to [%s]\n", wantedContext)
	},
}

func init() {
	authCmd.AddCommand(authUseContextCmd)
}
