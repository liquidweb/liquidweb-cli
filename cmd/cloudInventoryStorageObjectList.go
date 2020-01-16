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

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/types/api"
)

var cloudInventoryStorageObjectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Object Stores on your account",
	Long:  `List Object Stores on your account`,
	Run: func(cmd *cobra.Command, args []string) {
		methodArgs := instance.AllPaginatedResultsArgs{
			Method: "bleed/asset/list",
			MethodArgs: map[string]interface{}{
				"type": "SS.ObjectStore",
			},
			ResultsPerPage: 100,
		}
		results, err := lwCliInst.AllPaginatedResults(&methodArgs)
		if err != nil {
			lwCliInst.Die(err)
		}

		for _, item := range results.Items {

			itemUniqIdStr := cast.ToString(item["uniq_id"])

			var details apiTypes.CloudObjectStoreDetails
			apiArgs := map[string]interface{}{
				"uniq_id": itemUniqIdStr,
			}
			if err := lwCliInst.CallLwApiInto("bleed/storage/objectstore/details", apiArgs, &details); err != nil {
				lwCliInst.Die(err)
			}

			fmt.Printf(details.String())
		}
	},
}

func init() {
	cloudInventoryStorageObjectCmd.AddCommand(cloudInventoryStorageObjectListCmd)
}
