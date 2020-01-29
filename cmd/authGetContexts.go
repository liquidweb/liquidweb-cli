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

	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/types/cmd"
)

var authGetContextsCmd = &cobra.Command{
	Use:   "get-contexts",
	Short: "Display all auth contexts",
	Long: `Displays all configured auth contexts. 

If you've never setup any contexts, check "auth init".`,
	Run: func(cmd *cobra.Command, args []string) {
		contexts := lwCliInst.Viper.GetStringMap("liquidweb.api.contexts")

		for _, contextInter := range contexts {
			var context cmdTypes.AuthContext
			if err := instance.CastFieldTypes(contextInter, &context); err != nil {
				lwCliInst.Die(err)
			}

			fmt.Printf("Context: %s\n", context.ContextName)
			fmt.Printf("\tUsername: %s\n", context.Username)
			fmt.Printf("\tAPI URL: %s\n", context.Url)
			fmt.Printf("\tInsecure: %t\n", context.Insecure)
			fmt.Printf("\tTimeout: %d\n", context.Timeout)
		}

		currentContext := lwCliInst.Viper.GetString("liquidweb.api.current_context")
		fmt.Printf("Current context: [%s]\n", currentContext)
	},
}

func init() {
	authCmd.AddCommand(authGetContextsCmd)
}
