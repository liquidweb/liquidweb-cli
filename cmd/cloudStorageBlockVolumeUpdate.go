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
)

var cloudStorageBlockVolumeUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Set/Unset cross attach or rename a Cloud Block Storage Volume",
	Long: `Set/Unset cross attach or rename a Cloud Block Storage Volume.

Block storage offers a method to attach additional storage to Cloud Server.
Once attached, volumes appear as normal block devices, and can be used as such.
`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		nameFlag, _ := cmd.Flags().GetString("name")
		enableCrossAttachFlag, _ := cmd.Flags().GetBool("enable-cross-attach")
		disableCrossAttachFlag, _ := cmd.Flags().GetBool("disable-cross-attach")

		if enableCrossAttachFlag && disableCrossAttachFlag {
			lwCliInst.Die(fmt.Errorf("cant both enable and disab"))
		}

		apiArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
		}

		if enableCrossAttachFlag {
			apiArgs["cross_attach"] = true
		} else if disableCrossAttachFlag {
			apiArgs["cross_attach"] = false
		}

		if nameFlag != "" {
			apiArgs["domain"] = nameFlag
		}

		var details apiTypes.CloudBlockStorageVolumeDetails
		err := lwCliInst.CallLwApiInto("bleed/storage/block/volume/update",
			apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Updated Block Storage Volume %s\n%s", details.UniqId, details)
	},
}

func init() {
	cloudStorageBlockVolumeCmd.AddCommand(cloudStorageBlockVolumeUpdateCmd)

	cloudStorageBlockVolumeUpdateCmd.Flags().String("uniq_id", "",
		"uniq_id of Cloud Block Storage Volume")
	cloudStorageBlockVolumeUpdateCmd.Flags().String("name", "",
		"new name for the Cloud Block Storage Volume")
	cloudStorageBlockVolumeUpdateCmd.Flags().Bool("enable-cross-attach", false,
		"enable cross attach for Cloud Block Storage Volume")
	cloudStorageBlockVolumeUpdateCmd.Flags().Bool("disable-cross-attach", false,
		"disable cross attach for Cloud Block Storage Volume")

	cloudStorageBlockVolumeUpdateCmd.MarkFlagRequired("uniq_id")
}
