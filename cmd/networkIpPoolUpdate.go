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

var networkIpPoolUpdateCmdAddIpsFlag []string
var networkIpPoolUpdateCmdRemoveIpsFlag []string

var networkIpPoolUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing IP Pool",
	Long: `Update an existing IP Pool.

An IP Pool is a range of nonintersecting, reusable IP addresses reserved to
your account.`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		newIpsFlag, _ := cmd.Flags().GetInt64("new-ips")

		validateFields := map[interface{}]string{
			uniqIdFlag: "UniqId",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		if len(networkIpPoolUpdateCmdAddIpsFlag) == 0 && len(networkIpPoolUpdateCmdRemoveIpsFlag) == 0 &&
			newIpsFlag == -1 {
			lwCliInst.Die(fmt.Errorf(
				"at least one of --remove-ips --add-ips --new-ips flags must be given"))
		}

		apiArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
		}

		if len(networkIpPoolUpdateCmdAddIpsFlag) > 0 {
			apiArgs["add_ips"] = networkIpPoolUpdateCmdAddIpsFlag
		}
		if len(networkIpPoolUpdateCmdRemoveIpsFlag) > 0 {
			apiArgs["remove_ips"] = networkIpPoolUpdateCmdRemoveIpsFlag
		}
		if newIpsFlag != -1 {
			apiArgs["new_ips"] = newIpsFlag
		}

		var details apiTypes.NetworkIpPoolDetails
		err := lwCliInst.CallLwApiInto("bleed/network/pool/update", apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(details)
	},
}

func init() {
	networkIpPoolCmd.AddCommand(networkIpPoolUpdateCmd)

	networkIpPoolUpdateCmd.Flags().StringSliceVar(&networkIpPoolUpdateCmdRemoveIpsFlag, "remove-ips",
		[]string{}, "ips separated by ',' to remove from IP Pool")
	networkIpPoolUpdateCmd.Flags().StringSliceVar(&networkIpPoolUpdateCmdAddIpsFlag, "add-ips",
		[]string{}, "ips separated by ',' to add to IP Pool")
	networkIpPoolUpdateCmd.Flags().Int64("new-ips", -1, "amount of new IPs to assign to the IP Pool")
	networkIpPoolUpdateCmd.Flags().String("uniq_id", "", "uniq_id of IP Pool")

	networkIpPoolUpdateCmd.MarkFlagRequired("uniq_id")
}
