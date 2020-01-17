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

var networkIpPoolCreateCmdAddIpsFlag []string

var networkIpPoolCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an IP Pool",
	Long: `Create an IP Pool.

An IP Pool is a range of nonintersecting, reusable IP addresses reserved to
your account.`,
	Run: func(cmd *cobra.Command, args []string) {
		zoneFlag, _ := cmd.Flags().GetInt64("zone")
		newIpsFlag, _ := cmd.Flags().GetInt64("new-ips")

		//fmt.Printf("networkIpPoolCreateCmdAddIpsFlag: %+v\n", networkIpPoolCreateCmdAddIpsFlag)
		//fmt.Printf("%d %d\n", newIpsFlag, zoneFlag)

		if len(networkIpPoolCreateCmdAddIpsFlag) == 0 && newIpsFlag == -1 {
			lwCliInst.Die(fmt.Errorf("flags --new-ips --add-ips cannot both be empty"))
		}

		apiArgs := map[string]interface{}{
			"zone_id": zoneFlag,
		}

		if newIpsFlag != -1 {
			apiArgs["new_ips"] = newIpsFlag
		} else {
			apiArgs["add_ips"] = networkIpPoolCreateCmdAddIpsFlag
		}

		var details apiTypes.NetworkIpPoolDetails
		if err := lwCliInst.CallLwApiInto("bleed/network/pool/create", apiArgs, &details); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(details)
	},
}

func init() {
	networkIpPoolCmd.AddCommand(networkIpPoolCreateCmd)

	networkIpPoolCreateCmd.Flags().StringSliceVar(&networkIpPoolCreateCmdAddIpsFlag, "add-ips", []string{},
		"ips separated by ',' to add to created IP Pool")
	networkIpPoolCreateCmd.Flags().Int64("new-ips", -1, "amount of IPs to assign to the created IP Pool")
	networkIpPoolCreateCmd.Flags().Int64("zone", -1, "zone id to create the IP Pool in")

	networkIpPoolCreateCmd.MarkFlagRequired("zone")
}
