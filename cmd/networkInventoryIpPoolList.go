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
	//"os"

	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/types/api"
)

var networkInventoryIpPoolListCmd = &cobra.Command{
	Use:   "list",
	Short: "List IP Pools on your account",
	Long:  `List IP Pools on your account`,
	Run: func(cmd *cobra.Command, args []string) {
		methodArgs := instance.AllPaginatedResultsArgs{
			Method:         "bleed/network/pool/list",
			ResultsPerPage: 100,
		}
		results, err := lwCliInst.AllPaginatedResults(&methodArgs)
		if err != nil {
			lwCliInst.Die(err)
		}

		for _, item := range results.Items {

			var listEntry apiTypes.NetworkIpPoolListEntry
			if err := instance.CastFieldTypes(item, &listEntry); err != nil {
				lwCliInst.Die(err)
			}

			// now fetch details of Ip Pool id listEntry.Id
			var details apiTypes.NetworkIpPoolDetails
			if err := lwCliInst.CallLwApiInto("bleed/network/pool/details", map[string]interface{}{
				"id": listEntry.Id}, &details); err != nil {
				lwCliInst.Die(err)
			}

			fmt.Printf(details.String())
		}
	},
}

func init() {
	networkInventoryIpPoolCmd.AddCommand(networkInventoryIpPoolListCmd)
}
