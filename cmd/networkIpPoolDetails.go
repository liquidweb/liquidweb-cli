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

var networkIpPoolDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get details of an IP Pool",
	Long: `Get details of an IP Pool.

An IP Pool is a range of nonintersecting, reusable IP addresses reserved to
your account.`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")
		freeOnlyFlag, _ := cmd.Flags().GetBool("free-only")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag: "UniqId",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		apiArgs := map[string]interface{}{
			"uniq_id":   uniqIdFlag,
			"free_only": freeOnlyFlag,
		}

		var details apiTypes.NetworkIpPoolDetails
		if err := lwCliInst.CallLwApiInto("bleed/network/pool/details", apiArgs, &details); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(details)
	},
}

func init() {
	networkIpPoolCmd.AddCommand(networkIpPoolDetailsCmd)

	networkIpPoolDetailsCmd.Flags().String("uniq-id", "", "uniq-id of IP Pool")
	networkIpPoolDetailsCmd.Flags().Bool("free-only", false, "return only unassigned IPs in the IP Pool")

	networkIpPoolDetailsCmd.MarkFlagRequired("uniq-id")
}
