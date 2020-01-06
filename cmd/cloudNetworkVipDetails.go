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

var cloudNetworkVipDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get details of a Cloud Virtual IP (VIP)",
	Long: `Get details of a Cloud Virtual IP (VIP).

A Cloud Virtual IP (VIP) can be bound to multiple Cloud Servers. A common use case
for VIP is in High Availability setups where on failure the passive node becomes
active claiming the VIP.

When you create a Virtual IP (VIP) you will receive both a Public VIP and Private
VIP. The Public VIP can be configured on a Cloud Server just as a non-virtual, or
standard, IP would be configured. Connecting to a public service, such as HTTP or
FTP, on the Public VIP occurs just as it would on a standard IP.

The Private VIP can be configured on a Cloud Server’s private interface just as a
standard private IP would be configured. Connecting to a private service, such as
MySQL or Puppet, on the Private VIP also occurs just as it would on a standard
private IP.

So why use a VIP? When utilizing multiple servers, having a VIP is beneficial due
to its ability to “float” between Cloud Servers. This allows the VIP to remain
highly reachable in circumstances in which a non-virtual (or standard) IP may be
otherwise unreachable. It is possible to move both the Public VIP and Private VIP
between Cloud Servers, so long as they are in the same zone.

Use Cases for VIPs:

High Availability Databases (MySQL, Percona, MariaDB)
Non-DNS-based Service Migrations
High Availability Web Applications (in tandem with or in place of load balancer)

Common examples of high availability (HA) software often used in VIP setups:

Pacemaker
Heartbeat
`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")

		if uniqIdFlag == "" {
			lwCliInst.Die(fmt.Errorf("flag --uniq_id is required"))
		}

		apiArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
		}

		var details apiTypes.CloudNetworkVipDetails
		err := lwCliInst.CallLwApiInto("bleed/vip/details", apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		_printCloudNetworkVipDetails(&details)
	},
}

func _printCloudNetworkVipDetails(details *apiTypes.CloudNetworkVipDetails) {
	fmt.Printf("VIP Details:\n")
	fmt.Printf("\tActive: %d\n", details.Active)
	fmt.Printf("\tActive Status: %s\n", details.ActiveStatus)
	fmt.Printf("\tName: %s\n", details.Domain)
	fmt.Printf("\tUniqId: %s\n", details.UniqId)
	fmt.Printf("\tIP: %s\n", details.Ip)
	fmt.Printf("\tPrivate IP: %+v\n", details.PrivateIp)
}

func init() {
	cloudNetworkVipCmd.AddCommand(cloudNetworkVipDetailsCmd)
	cloudNetworkVipDetailsCmd.Flags().String("uniq_id", "", "uniq_id of VIP")
}
