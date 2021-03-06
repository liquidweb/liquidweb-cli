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

	"github.com/liquidweb/liquidweb-cli/flags/defaults"
)

var defaultFlagsNagOffCmd = &cobra.Command{
	Use:   "nags-off",
	Short: "Turn nags off",
	Long: `Turn nags off for unset default flags.

When a default flag is set (such as "zone") then any subcommand will use its
value in place if omitted. Default flags are auth context aware. For details
on auth contexts, see 'help auth'.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := defaults.NagsOff(); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Println("Nags turned off.")
	},
}

func init() {
	defaultFlagsCmd.AddCommand(defaultFlagsNagOffCmd)
}
