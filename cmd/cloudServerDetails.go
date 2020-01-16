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

var blockStorageVolumeList apiTypes.MergedPaginatedList
var fetchedBlockStorageVolumes bool

var cloudServerDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get details of a Cloud Server",
	Long: `Get details of a Cloud Server.

You can check this methods API documentation for what the returned fields mean:

https://cart.liquidweb.com/storm/api/docs/bleed/Storm/Server.html#method_details
`,
	Run: func(cmd *cobra.Command, args []string) {

		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		if uniqIdFlag == "" {
			lwCliInst.Die(fmt.Errorf("--uniq_id is a required flag"))
		}

		var details apiTypes.CloudServerDetails
		if err := lwCliInst.CallLwApiInto("bleed/storm/server/details",
			map[string]interface{}{"uniq_id": uniqIdFlag}, &details); err != nil {
			lwCliInst.Die(err)
		}

		if jsonFlag {
			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(details)
			if err != nil {
				lwCliInst.Die(err)
			}
			fmt.Printf(pretty)
			os.Exit(0)
		}

		_printExtendedCloudServerDetails(&details)
	},
}

func _printExtendedCloudServerDetails(details *apiTypes.CloudServerDetails) {
	fmt.Print(details)

	// private network
	var attachedDetails apiTypes.CloudNetworkPrivateIsAttachedResponse
	err := lwCliInst.CallLwApiInto("bleed/network/private/isattached", map[string]interface{}{
		"uniq_id": details.UniqId}, &attachedDetails)
	if err != nil {
		lwCliInst.Die(err)
	}
	fmt.Printf("\tPrivateNetwork: ")
	if attachedDetails.IsAttached {
		fmt.Printf("Attached\n")
	} else {
		fmt.Printf("None\n")
	}

	// block storage
	if !fetchedBlockStorageVolumes {
		methodArgs := instance.AllPaginatedResultsArgs{
			Method:         "bleed/storage/block/volume/list",
			ResultsPerPage: 100,
		}

		blockStorageVolumeList, err = lwCliInst.AllPaginatedResults(&methodArgs)
		if err != nil {
			lwCliInst.Die(err)
		}
		fetchedBlockStorageVolumes = true
	}
	fmt.Printf("\tBlock Storage Volumes:\n")
	for _, item := range blockStorageVolumeList.Items {
		var blockStorageDetails apiTypes.CloudBlockStorageVolumeDetails
		if err := instance.CastFieldTypes(item, &blockStorageDetails); err != nil {
			lwCliInst.Die(err)
		}

		for _, entry := range blockStorageDetails.AttachedTo {
			if entry.Resource == details.UniqId {
				fmt.Printf("\t\tVolume: %s\n", blockStorageDetails.Domain)
				fmt.Printf("\t\t\tUniqId: %s\n", blockStorageDetails.UniqId)
				fmt.Printf("\t\t\tSize: %d\n", blockStorageDetails.Size)
				fmt.Printf("\t\t\tStatus: %s\n", blockStorageDetails.Status)
				fmt.Printf("\t\t\tCross Attach Enabled: %t\n", blockStorageDetails.CrossAttach)
			}
		}
	}

}

func init() {
	cloudServerCmd.AddCommand(cloudServerDetailsCmd)

	cloudServerDetailsCmd.Flags().Bool("json", false, "output in json format")
	cloudServerDetailsCmd.Flags().String("uniq_id", "", "get details of this uniq_id")
}
