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

var cloudNetworkPublicRemoveCmdIpsFlag []string

var cloudNetworkPublicRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove Public IP(s) from a Cloud Server",
	Long: `Remove Public IP(s) from a Cloud Server.

Remove specific Public IP(s) from a Cloud Server. If the configure-ips flag is passed in,
the IP addresses given will also be automatically removed from the guest operating system.

If the configure-ips flag is not passed, the IP will be unassigned, and you will no longer
be able to route the IP. However the IP(s) will still be set in the guest operating system.
In this scenario, it will be up to the system administrator to remove the IP(s) from the
network configuration.

Note that you cannot remove the Cloud Servers primary ip with this command.`,
	Run: func(cmd *cobra.Command, args []string) {
		params := &instance.CloudNetworkPublicRemoveParams{}

		params.UniqId, _ = cmd.Flags().GetString("uniq-id")
		params.ConfigureIps, _ = cmd.Flags().GetBool("configure-ips")
		params.Ips = cloudNetworkPublicRemoveCmdIpsFlag

		status, err := lwCliInst.CloudNetworkPublicRemove(params)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(status)
	},
}

func init() {
	cloudNetworkPublicCmd.AddCommand(cloudNetworkPublicRemoveCmd)
	cloudNetworkPublicRemoveCmd.Flags().String("uniq-id", "", "uniq-id of the Cloud Server")
	cloudNetworkPublicRemoveCmd.Flags().Bool("configure-ips", false,
		"whether or not to automatically remove the IP address(es) in the server config")
	cloudNetworkPublicRemoveCmd.Flags().StringSliceVar(&cloudNetworkPublicRemoveCmdIpsFlag, "ips", []string{},
		"ips separated by ',' to remove from the Cloud Server")

	if err := cloudNetworkPublicRemoveCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
	if err := cloudNetworkPublicRemoveCmd.MarkFlagRequired("ips"); err != nil {
		lwCliInst.Die(err)
	}
}
