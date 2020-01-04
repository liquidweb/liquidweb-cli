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

var cloudServerBlockStorageOptimizedEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enables Cloud Block Storage Optimized on a Cloud Server",
	Long: `Enables Cloud Block Storage Optimized on a Cloud Server.

A Cloud Server uses more memory when using Cloud Block Storage volumes. Normal
Cloud Servers have enough pad for this to not be an issue, but Cloud Dedicated
purposely runs very tight. To work around this we can reduce the RAM allocated
to the Cloud Server to give more to the hypervisor. We call this Cloud Block
Storage Optimized.

Enabling Cloud Block Storage will cause your Cloud Server to reboot.`,

	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")

		if uniqIdFlag == "" {
			lwCliInst.Die(fmt.Errorf("flag --uniq_id is required"))
		}

		var optimized apiTypes.CloudServerIsBlockStorageOptimized
		if err := lwCliInst.CallLwApiInto("bleed/storm/server/issbsoptimized",
			map[string]interface{}{"uniq_id": uniqIdFlag}, &optimized); err != nil {
			lwCliInst.Die(err)
		}
		if optimized.IsOptimized {
			lwCliInst.Die(fmt.Errorf("Cloud Block Storage Optimized is already enabled on this Cloud Server"))
		}

		apiArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
			"value":   true,
		}

		var details apiTypes.CloudServerIsBlockStorageOptimizedSetResponse
		if err := lwCliInst.CallLwApiInto("bleed/storm/server/setsbsoptimized",
			apiArgs, &details); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Enabled Cloud Block Storage Optimized; %s\n", details.Updated)
	},
}

func init() {
	cloudServerBlockStorageOptimizedCmd.AddCommand(cloudServerBlockStorageOptimizedEnableCmd)
	cloudServerBlockStorageOptimizedEnableCmd.Flags().String("uniq_id", "", "uniq_id of Cloud Server")

}
