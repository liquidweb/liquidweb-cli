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

var cloudStorageBlockVolumeResizeCmd = &cobra.Command{
	Use:   "resize",
	Short: "Resize a Cloud Block Storage Volume",
	Long: `Resize a Cloud Block Storage Volume.

Block storage offers a method to attach additional storage to Cloud Server.
Once attached, volumes appear as normal block devices, and can be used as such.
`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		newSizeFlag, _ := cmd.Flags().GetInt64("new-size")

		if uniqIdFlag == "" || newSizeFlag == -1 {
			lwCliInst.Die(
				fmt.Errorf(
					"flags --uniq_id --new-size are required flags",
				))
		}

		apiArgs := map[string]interface{}{
			"uniq_id":  uniqIdFlag,
			"new_size": newSizeFlag,
		}

		var details apiTypes.CloudBlockStorageVolumeResize
		err := lwCliInst.CallLwApiInto("bleed/storage/block/volume/resize",
			apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Resized Block Storage Volume [%s] from size [%d] to [%d] GB\n",
			details.UniqId, details.OldSize, details.NewSize)
	},
}

func init() {
	cloudStorageBlockVolumeCmd.AddCommand(cloudStorageBlockVolumeResizeCmd)

	cloudStorageBlockVolumeResizeCmd.Flags().String("uniq_id", "",
		"uniq_id of Cloud Block Storage Volume")
	cloudStorageBlockVolumeResizeCmd.Flags().Int64("new-size", -1,
		"size (gb) to resize the volume to")
}
