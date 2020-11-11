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

var cloudNetworkPublicAddCmdPoolIpsFlag []string

var cloudNetworkPublicAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add Public IP(s) to a Cloud Server",
	Long: `Add Public IP(s) to a Cloud Server.

Add a number of IPs to an existing Cloud Server. If the configure-ips flag is
passed in, the IP addresses will be automatically configured within the guest
operating system.

If the configure-ips flag is not passed, the IP addresses will be assigned, and
routing will be allowed. However the IP(s) will not be automatically configured
in the guest operating system. In this scenario, it will be up to the system
administrator to add the IP(s) to the network configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		params := &instance.CloudNetworkPublicAddParams{}

		params.UniqId, _ = cmd.Flags().GetString("uniq-id")
		params.ConfigureIps, _ = cmd.Flags().GetBool("configure-ips")
		params.NewIps, _ = cmd.Flags().GetInt64("new-ips")
		params.PoolIps = cloudNetworkPublicAddCmdPoolIpsFlag

		status, err := lwCliInst.CloudNetworkPublicAdd(params)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(status)
	},
}

func init() {
	cloudNetworkPublicCmd.AddCommand(cloudNetworkPublicAddCmd)
	cloudNetworkPublicAddCmd.Flags().String("uniq-id", "", "uniq-id of the Cloud Server")
	cloudNetworkPublicAddCmd.Flags().Bool("configure-ips", false,
		"wheter or not to automatically configure the new IP address(es) in the server")
	cloudNetworkPublicAddCmd.Flags().Int64("new-ips", 0, "amount of new ips to (randomly) grab")
	cloudNetworkPublicAddCmd.Flags().StringSliceVar(&cloudNetworkPublicAddCmdPoolIpsFlag, "pool-ips", []string{},
		"ips from your IP Pool separated by ',' to assign to the Cloud Server")

	if err := cloudNetworkPublicAddCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
