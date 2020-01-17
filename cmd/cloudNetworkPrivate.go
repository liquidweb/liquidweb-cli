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

var cloudNetworkPrivateCmd = &cobra.Command{
	Use:   "private",
	Short: "Cloud Private Network specific operations",
	Long: `Cloud Private Network specific operations.

Private networking provides the option for several Cloud Servers to contact each other
via a network interface that is:

A) not publicly routable
B) not metered for bandwidth.

Applications that communicate internally will frequently use this for both security
and cost-savings.

For a full list of capabilities, please refer to the "Available Commands" section.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	cloudNetworkCmd.AddCommand(cloudNetworkPrivateCmd)
}
