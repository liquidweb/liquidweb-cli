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

var networkLoadBalancerAddServiceCmd = &cobra.Command{
	Use:   "add-service",
	Short: "Add a service to an existing Load Balancer",
	Long: `Add a service to an existing Load Balancer.

A service represents a service to load balance.`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")
		srcPortFlag, _ := cmd.Flags().GetInt("src-port")
		destPortFlag, _ := cmd.Flags().GetInt("dest-port")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag:   "UniqId",
			srcPortFlag:  "NetworkPort",
			destPortFlag: "NetworkPort",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		apiArgs := map[string]interface{}{
			"uniq_id":   uniqIdFlag,
			"src_port":  srcPortFlag,
			"dest_port": destPortFlag,
		}

		var details apiTypes.NetworkLoadBalancerDetails
		if err := lwCliInst.CallLwApiInto("bleed/network/loadbalancer/addservice", apiArgs,
			&details); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(details)
	},
}

func init() {
	networkLoadBalancerCmd.AddCommand(networkLoadBalancerAddServiceCmd)
	networkLoadBalancerAddServiceCmd.Flags().String("uniq-id", "", "uniq-id of Load Balancer")
	networkLoadBalancerAddServiceCmd.Flags().Int("src-port", -1, "source port")
	networkLoadBalancerAddServiceCmd.Flags().Int("dest-port", -1, "destination port")

	if err := networkLoadBalancerAddServiceCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
	if err := networkLoadBalancerAddServiceCmd.MarkFlagRequired("src-port"); err != nil {
		lwCliInst.Die(err)
	}
	if err := networkLoadBalancerAddServiceCmd.MarkFlagRequired("dest-port"); err != nil {
		lwCliInst.Die(err)
	}
}
