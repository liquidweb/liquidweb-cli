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
)

var authRemoveContextCmd = &cobra.Command{
	Use:   "remove-context",
	Short: "Remove a context from an existing configuration",
	Long: `Remove a context from an existing configuration.

Use this if you've already setup contexts with "auth init".`,
	Run: func(cmd *cobra.Command, args []string) {
		contextFlag, _ := cmd.Flags().GetString("context")
		if contextFlag == "" {
			lwCliInst.Die(fmt.Errorf("flag --context is required"))
		}

		if err := lwCliInst.RemoveContext(contextFlag); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Removed context [%s]\n", contextFlag)
	},
}

func init() {
	authCmd.AddCommand(authRemoveContextCmd)

	authRemoveContextCmd.Flags().String("context", "", "name of context to remove")
}
