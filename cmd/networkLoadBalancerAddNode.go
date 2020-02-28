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
	"github.com/liquidweb/liquidweb-cli/validate"
)

var networkLoadBalancerAddNodeCmd = &cobra.Command{
	Use:   "add-node",
	Short: "Add a node to an existing Load Balancer",
	Long: `Add a node (ip) to an existing Load Balancer.

A node is an ip address. You can only add ip addresses that are assigned to
your account.
`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")
		nodeFlag, _ := cmd.Flags().GetString("node")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag: "UniqId",
			nodeFlag:   "IP",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		apiArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
			"node":    nodeFlag,
		}

		var details apiTypes.NetworkLoadBalancerDetails
		if err := lwCliInst.CallLwApiInto("bleed/network/loadbalancer/addnode", apiArgs,
			&details); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(details)
	},
}

func init() {
	networkLoadBalancerCmd.AddCommand(networkLoadBalancerAddNodeCmd)
	networkLoadBalancerAddNodeCmd.Flags().String("uniq-id", "", "uniq-id of Load Balancer")
	networkLoadBalancerAddNodeCmd.Flags().String("node", "", "node (ip) to add to the Load Balancer")
	networkLoadBalancerAddNodeCmd.MarkFlagRequired("uniq-id")
	networkLoadBalancerAddNodeCmd.MarkFlagRequired("node")
}
