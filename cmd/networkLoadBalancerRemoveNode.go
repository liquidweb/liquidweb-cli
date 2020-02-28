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

var networkLoadBalancerRemoveNodeCmd = &cobra.Command{
	Use:   "remove-node",
	Short: "Remove a node from an existing Load Balancer",
	Long:  `Remove a node (ip) from an existing Load Balancer.`,
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
		if err := lwCliInst.CallLwApiInto("bleed/network/loadbalancer/removenode", apiArgs,
			&details); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(details)
	},
}

func init() {
	networkLoadBalancerCmd.AddCommand(networkLoadBalancerRemoveNodeCmd)
	networkLoadBalancerRemoveNodeCmd.Flags().String("uniq-id", "", "uniq-id of Load Balancer")
	networkLoadBalancerRemoveNodeCmd.Flags().String("node", "", "node (ip) to remove from the Load Balancer")
	networkLoadBalancerRemoveNodeCmd.MarkFlagRequired("uniq-id")
	networkLoadBalancerRemoveNodeCmd.MarkFlagRequired("node")
}
