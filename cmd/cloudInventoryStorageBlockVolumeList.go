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

var cloudInventoryStorageBlockVolumeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Cloud Block Storage volumes on your account",
	Long:  `List Cloud Block Storage volumes on your account`,
	Run: func(cmd *cobra.Command, args []string) {
		methodArgs := instance.AllPaginatedResultsArgs{
			Method:         "bleed/storage/block/volume/list",
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
			cnt := 1
			for _, item := range results.Items {

				var details apiTypes.CloudBlockStorageVolumeDetails
				if err := instance.CastFieldTypes(item, &details); err != nil {
					lwCliInst.Die(err)
				}

				fmt.Printf("%d.) %s", cnt, details)
				cnt++
			}
		}
	},
}

func init() {
	cloudInventoryStorageBlockVolumeCmd.AddCommand(cloudInventoryStorageBlockVolumeListCmd)

	cloudInventoryStorageBlockVolumeListCmd.Flags().Bool("json", false, "output in json format")
}
