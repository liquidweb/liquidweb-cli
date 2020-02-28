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

var networkLoadBalancerGetPossibleNodesCmd = &cobra.Command{
	Use:   "get-possible-nodes",
	Short: "Get possible Load Balancer nodes",
	Long: `Get possible nodes on the account.

When --region-id is passed, it will only list possible Load Balancer nodes
in that region.
`,
	Run: func(cmd *cobra.Command, args []string) {
		regionIdFlag, _ := cmd.Flags().GetInt("region-id")

		apiArgs := map[string]interface{}{}

		if regionIdFlag != -1 {
			validateFields := map[interface{}]interface{}{
				regionIdFlag: "PositiveInt",
			}
			if err := validate.Validate(validateFields); err != nil {
				lwCliInst.Die(err)
			}
			apiArgs["region"] = regionIdFlag
		}

		var details apiTypes.NetworkLoadBalancerPossibleNodes
		if err := lwCliInst.CallLwApiInto("bleed/network/loadbalancer/possiblenodes",
			apiArgs, &details); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(details)
	},
}

func init() {
	networkLoadBalancerCmd.AddCommand(networkLoadBalancerGetPossibleNodesCmd)
	networkLoadBalancerGetPossibleNodesCmd.Flags().Int("region-id", -1,
		"when passed only shows possible Load Balancer nodes in this region")
}
