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

	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/utils"
)

var cloudStorageBlockVolumeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Cloud Block Storage Volume",
	Long: `Create a Cloud Block Storage Volume.

Block storage offers a method to attach additional storage to Cloud Server.
Once attached, volumes appear as normal block devices, and can be used as such.
`,
	Run: func(cmd *cobra.Command, args []string) {
		nameFlag, _ := cmd.Flags().GetString("name")
		sizeFlag, _ := cmd.Flags().GetInt64("size")
		regionFlag, _ := cmd.Flags().GetInt64("region")
		zoneFlag, _ := cmd.Flags().GetInt64("zone")
		crossAttachFlag, _ := cmd.Flags().GetBool("cross-attach")
		attachFlag, _ := cmd.Flags().GetString("attach")

		if sizeFlag == -1 {
			lwCliInst.Die(fmt.Errorf("flag --size is required"))
		}

		apiArgs := map[string]interface{}{
			"domain":       nameFlag,
			"size":         sizeFlag,
			"cross_attach": crossAttachFlag,
		}

		if attachFlag != "" {
			apiArgs["attach"] = attachFlag
		}

		if regionFlag != -1 {
			apiArgs["region"] = regionFlag
		}
		if zoneFlag != -1 {
			apiArgs["zone"] = zoneFlag
		}

		var details apiTypes.CloudBlockStorageVolumeDetails
		err := lwCliInst.CallLwApiInto("bleed/storage/block/volume/create", apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Created Cloud Block Storage Volume\n%s", details.String())
	},
}

func init() {
	cloudStorageBlockVolumeCmd.AddCommand(cloudStorageBlockVolumeCreateCmd)

	cloudStorageBlockVolumeCreateCmd.Flags().Int64("size", -1, "size (gb) for Block Storage Volume")
	cloudStorageBlockVolumeCreateCmd.Flags().Int64("zone", -1, "Zone id for Block Storage Volume")
	cloudStorageBlockVolumeCreateCmd.Flags().Int64("region", -1, "Region id for Block Storage Volume")
	cloudStorageBlockVolumeCreateCmd.Flags().String("name", fmt.Sprintf("bsv-%s", utils.RandomString(5)),
		"Name for Block Storage volume")
	cloudStorageBlockVolumeCreateCmd.Flags().Bool("cross-attach", false, "Enable cross attach for Block Storage volume")
	cloudStorageBlockVolumeCreateCmd.Flags().String("attach", "", "uniq_id to attach created Block Storage volume to")
}
