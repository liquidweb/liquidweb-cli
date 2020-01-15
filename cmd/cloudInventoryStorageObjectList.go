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
	"github.com/liquidweb/liquidweb-cli/utils"
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

			_printCloudObjectStorageDetailsStruct(&details, itemUniqIdStr)
		}
	},
}

func _printCloudObjectStorageDetailsStruct(details *apiTypes.CloudObjectStoreDetails, uniqId string) {
	if uniqId == "" {
		uniqId = details.UniqId
	}

	var usage apiTypes.CloudObjectStoreDiskSpace
	if uniqId != "" {
		err := lwCliInst.CallLwApiInto("bleed/storage/objectstore/diskspace",
			map[string]interface{}{"uniq_id": uniqId}, &usage)
		if err != nil {
			lwCliInst.Die(err)
		}
	}

	utils.PrintTeal("Object Store UniqId: %s\n", uniqId)
	fmt.Printf("\tDisplay Name: %s\n", details.DisplayName)
	fmt.Printf("\tHost: %s\n", details.Host)
	fmt.Printf("\tUserId: %s\n", details.UserId)
	fmt.Printf("\tMax Buckets: %d\n", details.MaxBuckets)
	fmt.Printf("\tCaps:\n")
	for _, key := range details.Caps {
		fmt.Printf("\t\tPerm: %s\n", key.Perm)
		fmt.Printf("\t\tType: %s\n", key.Type)
	}
	fmt.Printf("\tKeys:\n")
	for _, key := range details.Keys {
		fmt.Printf("\t\tAccess Key: %s\n", key.AccessKey)
		fmt.Printf("\t\tSecret Key: %s\n", key.SecretKey)
		fmt.Printf("\t\tUser: %s\n", key.User)
	}
	if details.Accnt != 0 {
		fmt.Printf("\tAccount: %d\n", details.Accnt)
	}
	fmt.Printf("\tSuspended: %t\n", details.Suspended)

	if usage.Total != 0 {
		fmt.Printf("\tDiskspace total: %d\n", usage.Total)
		for _, bucket := range usage.Buckets {
			fmt.Printf("bucket: %+v\n", bucket)
		}
	}
}

func init() {
	cloudInventoryStorageObjectCmd.AddCommand(cloudInventoryStorageObjectListCmd)
}
