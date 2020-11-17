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

	"github.com/liquidweb/liquidweb-cli/instance"
)

var cloudNetworkPrivateDetachCmdUniqIdFlag []string

var cloudNetworkPrivateDetachCmd = &cobra.Command{
	Use:   "detach",
	Short: "Detach a Cloud Server from a Private Network",
	Long: `Detach a Cloud Server from a Private Network

Private networking provides the option for several Cloud Servers to contact each other
via a network interface that is:

A) not publicly routable
B) not metered for bandwidth.

Applications that communicate internally will frequently use this for both security
and cost-savings.
`,
	Run: func(cmd *cobra.Command, args []string) {
		params := &instance.CloudNetworkPrivateDetachParams{}

		params.UniqId = cloudNetworkPrivateDetachCmdUniqIdFlag

		status, err := lwCliInst.CloudNetworkPrivateDetach(params)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(status)
	},
}

func init() {
	cloudNetworkPrivateCmd.AddCommand(cloudNetworkPrivateDetachCmd)
	cloudNetworkPrivateDetachCmd.Flags().StringSliceVar(&cloudNetworkPrivateDetachCmdUniqIdFlag, "uniq-id",
		[]string{}, "uniq-ids separated by ',' of Cloud Servers to detach from private networking")
	if err := cloudNetworkPrivateDetachCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
