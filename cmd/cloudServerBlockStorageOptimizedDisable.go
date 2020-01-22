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

var cloudServerBlockStorageOptimizedDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disables Cloud Block Storage Optimized on a Cloud Server",
	Long: `Disables Cloud Block Storage Optimized on a Cloud Server.

A Cloud Server uses more memory when using Cloud Block Storage volumes. Normal
Cloud Servers have enough pad for this to not be an issue, but Cloud Dedicated
purposely runs very tight. To work around this we can reduce the RAM allocated
to the Cloud Server to give more to the hypervisor. We call this Cloud Block
Storage Optimized.

Disabling Cloud Block Storage will cause your Cloud Server to reboot.`,

	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")

		validateFields := map[string]interface{}{
			"UniqId": uniqIdFlag,
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(fmt.Errorf("flag validation failure: %s", err))
		}

		var optimized apiTypes.CloudServerIsBlockStorageOptimized
		if err := lwCliInst.CallLwApiInto("bleed/storm/server/issbsoptimized",
			map[string]interface{}{"uniq_id": uniqIdFlag}, &optimized); err != nil {
			lwCliInst.Die(err)
		}
		if !optimized.IsOptimized {
			lwCliInst.Die(fmt.Errorf("Cloud Block Storage Optimized is already not enabled on this Cloud Server"))
		}

		apiArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
			"value":   false,
		}

		var details apiTypes.CloudServerIsBlockStorageOptimized
		if err := lwCliInst.CallLwApiInto("bleed/storm/server/setsbsoptimized",
			apiArgs, &details); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Disabled Cloud Block Storage Optimized; %+v\n", details)
	},
}

func init() {
	cloudServerBlockStorageOptimizedCmd.AddCommand(cloudServerBlockStorageOptimizedDisableCmd)
	cloudServerBlockStorageOptimizedDisableCmd.Flags().String("uniq_id", "", "uniq_id of Cloud Server")

	cloudServerBlockStorageOptimizedDisableCmd.MarkFlagRequired("uniq_id")
}
