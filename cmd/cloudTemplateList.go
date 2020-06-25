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
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/instance"
)

var cloudTemplateListCmd = &cobra.Command{
	Use:   "list",
	Short: "Displays a list of cloud VPS templates",
	Long:  `Displays a list of cloud VPS templates.`,
	Run: func(cmd *cobra.Command, args []string) {
		zoneFlag, _ := cmd.Flags().GetInt("zone")
		filterOsFlag, _ := cmd.Flags().GetString("os")
		filterManageLevelFlag, _ := cmd.Flags().GetString("manage-level")

		templateList, err := lwCliInst.AllPaginatedResults(&instance.AllPaginatedResultsArgs{
			Method:         "bleed/storm/template/list",
			ResultsPerPage: 100,
		})
		if err != nil {
			lwCliInst.Die(err)
		}

		zonesResults, err := lwCliInst.LwCliApiClient.Call("bleed/network/zone/list", nil)
		if err != nil {
			lwCliInst.Die(err)
		}

		type ZoneInfo struct {
			Id         int
			Name       string
			RegionName string
		}
		zones := make(map[int]*ZoneInfo)
		zoneMap := zonesResults.(map[string]interface{})

		for _, item := range zoneMap["items"].([]interface{}) {
			z := item.(map[string]interface{})

			zone := &ZoneInfo{
				Id:         cast.ToInt(z["id"]),
				Name:       cast.ToString(z["name"]),
				RegionName: cast.ToString(z["region"].(map[string]interface{})["name"]),
			}
			zones[zone.Id] = zone
		}

		for _, template := range templateList.Items {
			if cast.ToBool(cast.ToInt(template["deprecated"])) {
				continue
			}

			if !strings.HasPrefix(strings.ToLower(cast.ToString(template["os"])), strings.ToLower(filterOsFlag)) {
				continue
			}

			if filterManageLevelFlag != "" {
				if strings.ToLower(filterManageLevelFlag) != strings.ToLower(cast.ToString(template["manage_level"])) {
					continue
				}
			}

			if zoneFlag != -1 {
				var skip bool = true

				for templateZoneStr, _ := range template["zone_availability"].(map[string]interface{}) {
					templateZone := cast.ToInt(templateZoneStr)
					if templateZone == zoneFlag {
						skip = false
					}
				}

				if skip {
					continue
				}
			}

			fmt.Println("name:", template["name"])
			fmt.Println("  description: ", template["description"])
			fmt.Print("  os: ", template["os"])
			fmt.Println(", manage-level:", template["manage_level"])
			fmt.Println("  Zone Availibility:")

			for templateZoneStr, _ := range template["zone_availability"].(map[string]interface{}) {
				templateZone := cast.ToInt(templateZoneStr)
				if z, ok := zones[templateZone]; ok {
					fmt.Printf("    %5d - %s - %s\n", z.Id, z.Name, z.RegionName)
				}

			}

			fmt.Println("")
		}
	},
}

func init() {
	cloudTemplateCmd.AddCommand(cloudTemplateListCmd)
	cloudTemplateListCmd.Flags().Int("zone", -1, "id of zone to filter by")
	cloudTemplateListCmd.Flags().String("os", "", "filter if os begins with string (i.e. linux, win)")
	cloudTemplateListCmd.Flags().String("manage-level", "", "filter list by management level")
}
