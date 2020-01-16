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

	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/spf13/cobra"
)

var cloudNetworkPublicAddCmdPoolIpsFlag []string

var cloudNetworkPublicAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add Public IP(s) to a Cloud Server",
	Long: `Add Public IP(s) to a Cloud Server.

Add a number of IPs to an existing Cloud Server. If the reboot flag is passed, the
server will be stopped, have the new IP addresses configured, and then started.

When the reboot flag is not passed, the IP will be assigned to the server, but it
will be up to the administrator to configure the IP address(es) within the server.
`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		rebootFlag, _ := cmd.Flags().GetBool("reboot")
		newIpsFlag, _ := cmd.Flags().GetInt64("new-ips")

		if newIpsFlag == 0 && len(cloudNetworkPublicAddCmdPoolIpsFlag) == 0 {
			lwCliInst.Die(fmt.Errorf("at least one of --new-ips --pool-ips must be given"))
		}

		apiArgs := map[string]interface{}{
			"reboot":  rebootFlag,
			"uniq_id": uniqIdFlag,
		}
		if newIpsFlag != 0 {
			apiArgs["ip_count"] = newIpsFlag
		}
		if len(cloudNetworkPublicAddCmdPoolIpsFlag) != 0 {
			apiArgs["pool_ips"] = cloudNetworkPublicAddCmdPoolIpsFlag
		}

		var details apiTypes.NetworkIpAdd
		err := lwCliInst.CallLwApiInto("bleed/network/ip/add", apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Adding [%s] to Cloud Server\n", details.Adding)
	},
}

func init() {
	cloudNetworkPublicCmd.AddCommand(cloudNetworkPublicAddCmd)
	cloudNetworkPublicAddCmd.Flags().String("uniq_id", "", "uniq_id of the Cloud Server")
	cloudNetworkPublicAddCmd.Flags().Bool("reboot", false,
		"wheter or not to automatically configure the new IP address(es) in the server (requires reboot)")
	cloudNetworkPublicAddCmd.Flags().Int64("new-ips", 0, "amount of new ips to (randomly) grab")
	cloudNetworkPublicAddCmd.Flags().StringSliceVar(&cloudNetworkPublicAddCmdPoolIpsFlag, "pool-ips", []string{},
		"ips from your IP Pool separated by ',' to assign to the Cloud Server")

	cloudNetworkPublicAddCmd.MarkFlagRequired("uniq_id")
}
