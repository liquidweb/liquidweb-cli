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
	"os"

	"github.com/spf13/cobra"
)

var networkIpPoolCmd = &cobra.Command{
	Use:   "ip-pool",
	Short: "IP Pool specific operations",
	Long: `IP Pool specific operations.

Use an IP Pool to reserve IP addresses even when they aren't currently assigned
to a server.

For a full list of capabilities, please refer to the "Available Commands" section.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			lwCliInst.Die(err)
		}
		os.Exit(1)
	},
}

func init() {
	networkCmd.AddCommand(networkIpPoolCmd)
}
