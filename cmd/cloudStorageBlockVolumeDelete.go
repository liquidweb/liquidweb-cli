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
	"github.com/liquidweb/liquidweb-cli/validate"
)

var cloudStorageBlockVolumeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Cloud Block Storage Volume",
	Long: `Delete a Cloud Block Storage Volume.

Block storage offers a method to attach additional storage to Cloud Server.
Once attached, volumes appear as normal block devices, and can be used as such.
`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")

		validateFields := map[string]interface{}{
			"UniqId": uniqIdFlag,
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(fmt.Errorf("flag validation failure: %s", err))
		}

		apiArgs := map[string]interface{}{"uniq_id": uniqIdFlag}
		var details apiTypes.CloudBlockStorageVolumeDelete
		err := lwCliInst.CallLwApiInto("bleed/storage/block/volume/delete", apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Deleted Cloud Block Storage Volume: %s\n", details.Deleted)
	},
}

func init() {
	cloudStorageBlockVolumeCmd.AddCommand(cloudStorageBlockVolumeDeleteCmd)

	cloudStorageBlockVolumeDeleteCmd.Flags().String("uniq_id", "", "uniq_id of Cloud Block Storage volume")

	cloudStorageBlockVolumeDeleteCmd.MarkFlagRequired("uniq_id")
}
