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
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/instance"
)

var cloudServerOptionsCmd = &cobra.Command{
	Use:   "options",
	Short: "Fetches available regions, zones, and templates",
	Long: `Fetches available regions, zones, and templates.

Use this when you need to get a list of available:

*) Regions/Zones
*) Templates (cloud server images provided by LiquidWeb)
*) config-id's (a config-id represents a type of server, for example, config-id 1234 might
   represent a configuration with the following hardware specifications:
     *) 16GB RAM
     *) 300GB Disk
     *) 8 vcpus

Be sure to take a look at the flags section for specific flags to pass.`,
	Run: func(cmd *cobra.Command, args []string) {
		configsFlag, _ := cmd.Flags().GetBool("configs")
		configCategoryFlag, _ := cmd.Flags().GetString("config-category")
		zonesFlag, _ := cmd.Flags().GetBool("zones")
		templatesFlag, _ := cmd.Flags().GetBool("templates")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		if !configsFlag && !zonesFlag && !templatesFlag && !jsonFlag {
			if err := cmd.Help(); err != nil {
				lwCliInst.Die(err)
			}
			os.Exit(1)
		}

		// buildout regionsWithZoneInfo before proceeding
		regionsWithZoneInfo := map[int][]map[string]interface{}{}

		cfgListArgs := instance.AllPaginatedResultsArgs{
			Method: "bleed/storm/config/list",
			MethodArgs: map[string]interface{}{
				"category":  configCategoryFlag,
				"available": 1,
			},
			ResultsPerPage: 100,
		}
		mergedConfigs, err := lwCliInst.AllPaginatedResults(&cfgListArgs)
		if err != nil {
			lwCliInst.Die(err)
		}

		mergedTemplates, err := lwCliInst.AllPaginatedResults(&instance.AllPaginatedResultsArgs{
			Method:         "bleed/storm/template/list",
			ResultsPerPage: 100,
		})
		if err != nil {
			lwCliInst.Die(err)
		}

		// add config and zone data to regionsWithZoneInfo
		zones := map[int]int{}
		configIdsByZone := map[int][]map[string]interface{}{}
		for _, item := range mergedConfigs.Items {
			if _, zoneAvailExists := item["zone_availability"]; !zoneAvailExists {
				continue
			}
			for key, value := range item["zone_availability"].(map[string]interface{}) {
				if cast.ToInt(value) == 1 {
					zone := cast.ToInt(key)
					if _, exists := zones[zone]; !exists {
						zones[zone] = 1
					}

					configDescription := cast.ToString(item["description"])
					if strings.Contains(configDescription, "VPS - KA") {
						continue
					} else if strings.Contains(configDescription, "Hybrid Dedicated") {
						continue
					}

					configData := map[string]interface{}{
						"config_id":   cast.ToInt(item["id"]),
						"active":      cast.ToInt(item["active"]),
						"available":   cast.ToInt(item["available"]),
						"category":    cast.ToString(item["category"]),
						"description": configDescription,
						"featured":    cast.ToInt(item["featured"]),
					}

					if configData["category"] == "bare-metal" || configData["category"] == "bare-metal-r" {
						configData["disk_count"] = cast.ToInt(item["disk_count"])
						configData["disk_total"] = cast.ToInt(item["disk_total"])
						configData["disk_type"] = cast.ToString(item["disk_type"])
						configData["ram_available"] = cast.ToInt(item["ram_available"])
						configData["ram_total"] = cast.ToInt(item["ram_total"])
						configData["cpu_cores"] = cast.ToInt(item["cpu_cores"])
						configData["cpu_count"] = cast.ToInt(item["cpu_count"])
						configData["cpu_hyperthreading"] = cast.ToInt(item["cpu_hyperthreading"])
						configData["cpu_model"] = cast.ToString(item["cpu_model"])
						configData["cpu_speed"] = cast.ToInt(item["cpu_speed"])
					} else {
						configData["disk"] = cast.ToInt(item["disk"])
						configData["memory"] = cast.ToInt(item["memory"])
						configData["vcpu"] = cast.ToInt(item["vcpu"])
					}

					configIdsByZone[zone] = append(configIdsByZone[zone], configData)
				}
			}
		}

		// determine template availability by zone
		templatesByZone := map[int][]map[string]interface{}{}
		for _, template := range mergedTemplates.Items {
			for templateZoneStr, _ := range template["zone_availability"].(map[string]interface{}) {
				if cast.ToInt(template["deprecated"]) != 0 {
					continue
				}

				templateZone := cast.ToInt(templateZoneStr)
				templateData := map[string]interface{}{
					"name":         template["name"],
					"os":           template["os"],
					"id":           cast.ToInt(template["id"]),
					"manage_level": template["manage_level"],
					"description":  template["description"],
				}
				templatesByZone[templateZone] = append(templatesByZone[templateZone], templateData)
			}
		}

		// add final region, config, and template data to regionsWithZoneInfo
		for zone, _ := range zones {
			zoneDetails, err := lwCliInst.LwCliApiClient.Call("bleed/network/zone/details", map[string]interface{}{"id": zone})
			if err != nil {
				lwCliInst.Die(err)
			}

			if zoneDetails.(map[string]interface{})["status"] != "Open" {
				continue
			}

			regionId := cast.ToInt(zoneDetails.(map[string]interface{})["region"].(map[string]interface{})["id"])
			regionsWithZoneInfo[regionId] = append(regionsWithZoneInfo[regionId], map[string]interface{}{
				"regionName": zoneDetails.(map[string]interface{})["region"].(map[string]interface{})["name"],
				"regionId":   regionId,
				"zoneId":     zone,
				"zoneName":   zoneDetails.(map[string]interface{})["name"],
				"configIds":  configIdsByZone[zone],
				"templates":  templatesByZone[zone],
			})
		}

		if jsonFlag {
			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(regionsWithZoneInfo)
			if err != nil {
				lwCliInst.Die(err)
			}
			fmt.Printf(pretty)
			os.Exit(0)
		}

		for region, dataSlice := range regionsWithZoneInfo {

			for _, info := range dataSlice {
				fmt.Printf("Working region [%s (id: %d)] zone [%s (id: %d)]\n", info["regionName"], region, info["zoneName"], info["zoneId"])

				// display available configs
				if configsFlag {
					fmt.Printf("  configs:\n")
					for _, cfgInfo := range info["configIds"].([]map[string]interface{}) {
						fmt.Printf("    config-id: %d\n", cfgInfo["config_id"])
						fmt.Printf("      description: %s\n", cfgInfo["description"])
						fmt.Printf("      active: %d\n", cfgInfo["active"])
						fmt.Printf("      available: %d\n", cfgInfo["available"])
						fmt.Printf("      category: %s\n", cfgInfo["category"])
						if cfgInfo["category"] == "bare-metal" || cfgInfo["category"] == "bare-metal-r" {
							fmt.Printf("      disk_count: %d\n", cfgInfo["disk_count"])
							fmt.Printf("      disk_total: %d\n", cfgInfo["disk_total"])
							fmt.Printf("      disk_type: %s\n", cfgInfo["disk_type"])
							fmt.Printf("      ram_available: %d\n", cfgInfo["ram_available"])
							fmt.Printf("      ram_total: %d\n", cfgInfo["ram_total"])
							fmt.Printf("      cpu_cores: %d\n", cfgInfo["cpu_cores"])
							fmt.Printf("      cpu_count: %d\n", cfgInfo["cpu_count"])
							fmt.Printf("      cpu_hyperthreading: %d\n", cfgInfo["cpu_hyperthreading"])
							fmt.Printf("      cpu_model: %s\n", cfgInfo["cpu_model"])
							fmt.Printf("      cpu_speed: %d\n", cfgInfo["cpu_speed"])
						} else {
							fmt.Printf("      disk: %d\n", cfgInfo["disk"])
							fmt.Printf("      memory: %d\n", cfgInfo["memory"])
							fmt.Printf("      vcpu: %d\n", cfgInfo["vcpu"])
						}
					}
				}

				// display available regions and their zones
				if zonesFlag {
					fmt.Printf("  zones:\n")
					fmt.Printf("    Region: %s (id: %d) Zone: %s (id: %d)\n", info["regionName"], info["regionId"], info["zoneName"], info["zoneId"])
				}

				// display available templates
				if templatesFlag {
					fmt.Printf("  templates:\n")
					for _, templateInfo := range info["templates"].([]map[string]interface{}) {
						fmt.Printf("    name: %s\n", templateInfo["name"])
						fmt.Printf("      description: %s\n", templateInfo["description"])
						fmt.Printf("      manage_level: %s\n", templateInfo["manage_level"])
						fmt.Printf("      os: %s\n", templateInfo["os"])
						fmt.Printf("      id: %d\n", templateInfo["id"])
					}
				}
			}
		}
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerOptionsCmd)

	cloudServerOptionsCmd.Flags().Bool("json", false, "return data in json format. All template, config, zone, and region data will be returned with this option.")
	cloudServerOptionsCmd.Flags().Bool("configs", false, "fetch a list of available configs (config-id)")
	cloudServerOptionsCmd.Flags().String("config-category", "all", "valid options for category are storm, ssd, bare-metal and all. Only relevent when --configs is passed.")
	cloudServerOptionsCmd.Flags().Bool("zones", false, "fetch a list of available regions and their available zones")
	cloudServerOptionsCmd.Flags().Bool("templates", false, "fetch a list of available templates (cloud server images provided by LiquidWeb)")
}
