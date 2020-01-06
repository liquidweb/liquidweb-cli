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
	"os"

	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/types/api"
)

var cloudInventoryNetworkVipListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VIPs on your account",
	Long:  `List all VIPs on your account`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		methodArgs := instance.AllPaginatedResultsArgs{
			Method: "bleed/asset/list",
			MethodArgs: map[string]interface{}{
				"type":     "SS.VIP",
				"alsowith": []string{"zone"},
			},
			ResultsPerPage: 100,
		}
		results, err := lwCliInst.AllPaginatedResults(&methodArgs)
		if err != nil {
			lwCliInst.Die(err)
		}

		if jsonFlag {
			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(results)
			if err != nil {
				lwCliInst.Die(err)
			}
			fmt.Printf(pretty)
			os.Exit(0)
		}

		for _, item := range results.Items {
			var details apiTypes.CloudNetworkVipAssetListAlsoWithZoneResponse
			if err := instance.CastFieldTypes(item, &details); err != nil {
				lwCliInst.Die(err)
			}

			fmt.Printf("VIP Details:\n")
			fmt.Printf("\tActive: %d\n", details.Active)
			fmt.Printf("\tName: %s\n", details.Domain)
			fmt.Printf("\tUniqId: %s\n", details.UniqId)
			fmt.Printf("\tIP: %s\n", details.Ip)
			fmt.Printf("\tRegion %s (id %d) Zone %s (id %d)\n", details.Zone.Region.Name,
				details.Zone.Region.Id, details.Zone.Name, details.Zone.Id)
		}
	},
}

func init() {
	cloudInventoryNetworkVipCmd.AddCommand(cloudInventoryNetworkVipListCmd)

	cloudInventoryNetworkVipListCmd.Flags().Bool("json", false, "output in json format")
}
