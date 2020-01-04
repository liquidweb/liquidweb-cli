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

var cloudInventoryBackupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Cloud Backups on your account",
	Long:  `List Cloud Backups on your account`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")

		methodArgs := instance.AllPaginatedResultsArgs{
			Method:         "bleed/storm/backup/list",
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
			var details apiTypes.CloudBackupDetails
			if err := instance.CastFieldTypes(item, &details); err != nil {
				lwCliInst.Die(err)
			}

			if uniqIdFlag != "" {
				if details.UniqId != uniqIdFlag {
					continue
				}
			}
			_printCloudBackupDetailsFromDetailsStruct(&details)
		}
	},
}

func _printCloudBackupDetailsFromDetailsStruct(details *apiTypes.CloudBackupDetails) {
	fmt.Printf("Cloud Backup Id: %d\n", details.Id)
	fmt.Printf("\tUniqId: %s\n", details.UniqId)
	fmt.Printf("\tName: %s\n", details.Name)
	fmt.Printf("\tsize: %.2f\n", details.Size)
	fmt.Printf("\tTemplate: %s\n", details.Template)
	fmt.Printf("\tTime Taken: %s\n", details.TimeTaken)
	fmt.Printf("\tHypervisor: %s\n", details.HvType)
	fmt.Printf("\tFeatures: %+v\n", details.Features)
	fmt.Printf("\tAccount: %d\n", details.Accnt)
}

func init() {
	cloudInventoryBackupCmd.AddCommand(cloudInventoryBackupListCmd)

	cloudInventoryBackupListCmd.Flags().Bool("json", false, "output in json format")
	cloudInventoryBackupListCmd.Flags().String("uniq_id", "", "only fetch backups made from this uniq_id")
}
