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

var cloudStorageBlockVolumeDetachCmd = &cobra.Command{
	Use:   "detach",
	Short: "Detach a Cloud Block Storage Volume from a Cloud Server",
	Long: `Detach a Cloud Block Storage Volume from a Cloud Server.

Block storage offers a method to attach additional storage to Cloud Server.
Once attached, volumes appear as normal block devices, and can be used as such.
`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")
		detachFromFlag, _ := cmd.Flags().GetString("detach-from")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag:     "UniqId",
			detachFromFlag: "UniqId",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		apiArgs := map[string]interface{}{
			"uniq_id":     uniqIdFlag,
			"detach_from": detachFromFlag,
		}

		var details apiTypes.CloudBlockStorageVolumeDetach
		err := lwCliInst.CallLwApiInto("bleed/storage/block/volume/detach",
			apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Detached Block Storage Volume %s from Cloud Server %s\n",
			details.Detached, details.DetachedFrom)
	},
}

func init() {
	cloudStorageBlockVolumeCmd.AddCommand(cloudStorageBlockVolumeDetachCmd)

	cloudStorageBlockVolumeDetachCmd.Flags().String("uniq-id", "",
		"uniq-id of Cloud Block Storage Volume")
	cloudStorageBlockVolumeDetachCmd.Flags().String("detach-from", "",
		"uniq-id of Cloud Server to detach from")

	if err := cloudStorageBlockVolumeDetachCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
	if err := cloudStorageBlockVolumeDetachCmd.MarkFlagRequired("detach-from"); err != nil {
		lwCliInst.Die(err)
	}
}
