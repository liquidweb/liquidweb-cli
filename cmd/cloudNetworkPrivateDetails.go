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

	"github.com/liquidweb/liquidweb-cli/types/api"
)

var cloudNetworkPrivateDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get Private Network details for a Cloud Server",
	Long: `Get Private Network details for a Cloud Server

Private networking provides the option for several Cloud Servers to contact each other
via a network interface that is:

A) not publicly routable
B) not metered for bandwidth.

Applications that communicate internally will frequently use this for both security
and cost-savings.
`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")

		apiArgs := map[string]interface{}{"uniq_id": uniqIdFlag}

		var details apiTypes.CloudNetworkPrivateGetIpResponse
		err := lwCliInst.CallLwApiInto("bleed/network/private/getip", apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		if details.Ip == "" {
			fmt.Printf("Cloud Server is not attached to a Private Network\n")
		} else {
			fmt.Printf("Cloud Server is attached to a Private Network\n")
			fmt.Printf("\tIP: %s\n", details.Ip)
			fmt.Printf("\tLegacy: %t\n", details.Legacy)
		}
	},
}

func init() {
	cloudNetworkPrivateCmd.AddCommand(cloudNetworkPrivateDetailsCmd)
	cloudNetworkPrivateDetailsCmd.Flags().String("uniq_id", "", "uniq_id of the Cloud Server")
	cloudNetworkPrivateDetailsCmd.MarkFlagRequired("uniq_id")
}
