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

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/types/api"
)

var cloudNetworkPublicListCmdPoolIpsFlag []string

var cloudNetworkPublicListCmd = &cobra.Command{
	Use:   "list",
	Short: "List a Cloud Servers Public IP(s)",
	Long:  `List a Cloud Servers Public IP(s).`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")

		methodArgs := instance.AllPaginatedResultsArgs{
			Method:         "bleed/network/ip/list",
			ResultsPerPage: 100,
			MethodArgs: map[string]interface{}{
				"uniq_id":    uniqIdFlag,
				"expand_ips": 1,
			},
		}
		results, err := lwCliInst.AllPaginatedResults(&methodArgs)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("IP Assignments for %s:\n\n", uniqIdFlag)

		for c, item := range results.Items {
			var details apiTypes.NetworkAssignmentListEntry
			if err := instance.CastFieldTypes(item, &details); err != nil {
				lwCliInst.Die(err)
			}

			// first ip is always primary
			if c == 0 {
				fmt.Println("Primary IP:")
			} else {
				fmt.Println("Secondary IP:")
			}
			fmt.Print(details)
		}
	},
}

func init() {
	cloudNetworkPublicCmd.AddCommand(cloudNetworkPublicListCmd)
	cloudNetworkPublicListCmd.Flags().String("uniq-id", "", "uniq-id of the Cloud Server")
	if err := cloudNetworkPublicListCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
