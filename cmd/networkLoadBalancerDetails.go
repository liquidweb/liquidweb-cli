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

var networkLoadBalancerDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get details of a Load Balancer",
	Long:  `Get details of a Load Balancer.`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag: "UniqId",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		apiArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
		}

		var details apiTypes.NetworkLoadBalancerDetails
		if err := lwCliInst.CallLwApiInto("bleed/network/loadbalancer/details", apiArgs,
			&details); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(details)
	},
}

func init() {
	networkLoadBalancerCmd.AddCommand(networkLoadBalancerDetailsCmd)
	networkLoadBalancerDetailsCmd.Flags().String("uniq-id", "", "uniq-id of Load Balancer")
	if err := networkLoadBalancerDetailsCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
