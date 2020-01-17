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

var cloudStorageBlockVolumeAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a Cloud Block Storage Volume to a Cloud Server",
	Long: `Attach a Cloud Block Storage Volume to a Cloud Server.

Block storage offers a method to attach additional storage to Cloud Server.
Once attached, volumes appear as normal block devices, and can be used as such.
`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		attachToFlag, _ := cmd.Flags().GetString("attach-to")

		apiArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
			"to":      attachToFlag,
		}

		var details apiTypes.CloudBlockStorageVolumeAttach
		err := lwCliInst.CallLwApiInto("bleed/storage/block/volume/attach",
			apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Attached Block Storage Volume %s to Cloud Server %s\n",
			details.Attached, details.To)
	},
}

func init() {
	cloudStorageBlockVolumeCmd.AddCommand(cloudStorageBlockVolumeAttachCmd)

	cloudStorageBlockVolumeAttachCmd.Flags().String("uniq_id", "",
		"uniq_id of Cloud Block Storage Volume")
	cloudStorageBlockVolumeAttachCmd.Flags().String("attach-to", "",
		"uniq_id of Cloud Server to attach to")

	cloudStorageBlockVolumeAttachCmd.MarkFlagRequired("uniq_id")
	cloudStorageBlockVolumeAttachCmd.MarkFlagRequired("attach-to")
}
