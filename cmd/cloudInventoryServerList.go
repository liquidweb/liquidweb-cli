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
)

var cloudInventoryServerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Cloud Servers on your account",
	Long:  `List Cloud Servers on your account`,
	Run: func(cmd *cobra.Command, args []string) {
		methodArgs := instance.AllPaginatedResultsArgs{
			Method:         "bleed/storm/server/list",
			ResultsPerPage: 100,
		}
		results, err := lwCliInst.AllPaginatedResults(&methodArgs)
		if err != nil {
			lwCliInst.Die(err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {

			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(results)
			if err != nil {
				lwCliInst.Die(err)
			}
			fmt.Printf(pretty)
		} else {
			serverCnt := 1
			for _, item := range results.Items {
				fmt.Printf("%d.) domain: %s / uniq_id: %s\n", serverCnt, item["domain"], item["uniq_id"])

				if zoneInter, zoneExists := item["zone"]; zoneExists {
					zone := zoneInter.(map[string]interface{})
					if regionInter, regionExists := zone["region"]; regionExists {
						region := regionInter.(map[string]interface{})
						fmt.Printf("\tzone: %+v - %+v\n", region["name"], zone["name"])
						fmt.Printf("\t\tzone_id: %+v\n", zone["id"])
						fmt.Printf("\t\tregion_id: %+v\n", region["id"])
					}
				}

				fields := []string{
					"create_date",
					"config_description",
					"config_id",
					"template",
					"template_description",
					"manage_level",
					"type",
					"backup_plan",
					"backup_quota",
					"backup_size",
					"bandwidth_quota",
					"diskspace",
					"memory",
					"vcpu",
					"ip",
					"ip_count",
				}
				for _, field := range fields {
					fmt.Printf("\t%s: %+v\n", field, item[field])
				}

				serverCnt++
			}
		}
	},
}

func init() {
	cloudInventoryServerCmd.AddCommand(cloudInventoryServerListCmd)

	cloudInventoryServerListCmd.Flags().Bool("json", false, "output in json format")
}
