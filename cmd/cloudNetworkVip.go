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

var cloudNetworkVipCmd = &cobra.Command{
	Use:   "vip",
	Short: "Cloud Virtual IP (VIP) specific operations",
	Long: `Cloud Virtual IP (VIP) specific operations.

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

For a full list of capabilities, please refer to the "Available Commands" section.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			lwCliInst.Die(err)
		}
		os.Exit(1)
	},
}

func init() {
	cloudNetworkCmd.AddCommand(cloudNetworkVipCmd)
}
