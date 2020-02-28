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

	apiTypes "github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/validate"
)

var cloudNetworkPublicRemoveCmdIpsFlag []string

var cloudNetworkPublicRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove Public IP(s) from a Cloud Server",
	Long: `Remove Public IP(s) from a Cloud Server.

Remove specific Public IP(s) from a Cloud Server. If the reboot flag is passed in, the machine
will be stopped, have the old IP addresses removed, and then started.

If the reboot flag is not passed, the IP will be unassigned, and you will no longer be able
to route the IP. However the machine will not be shutdown to remove it from its network
configuration. It will be up to the administrator to remove the IP from the servers network
configuration.

Note that you cannot remove the Cloud Servers primary ip with this command.`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")
		rebootFlag, _ := cmd.Flags().GetBool("reboot")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag: "UniqId",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		apiArgs := map[string]interface{}{
			"reboot":  rebootFlag,
			"uniq_id": uniqIdFlag,
		}

		for _, ip := range cloudNetworkPublicRemoveCmdIpsFlag {
			validateFields := map[interface{}]interface{}{
				ip: "IP",
			}
			if err := validate.Validate(validateFields); err != nil {
				fmt.Printf("%s ... skipping\n", err)
				continue
			}

			var details apiTypes.NetworkIpRemove
			apiArgs["ip"] = ip
			err := lwCliInst.CallLwApiInto("bleed/network/ip/remove", apiArgs, &details)
			if err != nil {
				lwCliInst.Die(err)
			}

			fmt.Printf("Removing [%s] from Cloud Server\n", details.Removing)
		}
	},
}

func init() {
	cloudNetworkPublicCmd.AddCommand(cloudNetworkPublicRemoveCmd)
	cloudNetworkPublicRemoveCmd.Flags().String("uniq-id", "", "uniq-id of the Cloud Server")
	cloudNetworkPublicRemoveCmd.Flags().Bool("reboot", false,
		"whether or not to automatically remove the IP address(es) in the server config (requires reboot)")
	cloudNetworkPublicRemoveCmd.Flags().StringSliceVar(&cloudNetworkPublicRemoveCmdIpsFlag, "ips", []string{},
		"ips separated by ',' to remove from the Cloud Server")

	cloudNetworkPublicRemoveCmd.MarkFlagRequired("uniq-id")
	cloudNetworkPublicRemoveCmd.MarkFlagRequired("ips")
}
