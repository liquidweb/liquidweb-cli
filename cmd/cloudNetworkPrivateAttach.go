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

var cloudNetworkPrivateAttachCmdUniqIdFlag []string

var cloudNetworkPrivateAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a Cloud Server to a Private Network",
	Long: `Attach a Cloud Server to a Private Network

Private networking provides the option for several Cloud Servers to contact each other
via a network interface that is:

A) not publicly routable
B) not metered for bandwidth.

Applications that communicate internally will frequently use this for both security
and cost-savings.
`,
	Run: func(cmd *cobra.Command, args []string) {
		params := &instance.CloudNetworkPrivateAttachParams{}

		params.UniqId = cloudNetworkPrivateAttachCmdUniqIdFlag

		status, err := lwCliInst.CloudNetworkPrivateAttach(params)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(status)
	},
}

func init() {
	cloudNetworkPrivateCmd.AddCommand(cloudNetworkPrivateAttachCmd)
	cloudNetworkPrivateAttachCmd.Flags().StringSliceVar(&cloudNetworkPrivateAttachCmdUniqIdFlag, "uniq-id",
		[]string{}, "uniq-ids separated by ',' of Cloud Servers to attach to private networking")
	if err := cloudNetworkPrivateAttachCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
