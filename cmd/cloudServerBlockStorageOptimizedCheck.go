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

var cloudServerBlockStorageOptimizedCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if a Cloud Server is Cloud Block Storage Optimized",
	Long: `Check if a Cloud Server is Cloud Block Storage Optimized.

A Cloud Server uses more memory when using Cloud Block Storage volumes. Normal
Cloud Servers have enough pad for this to not be an issue, but Cloud Dedicated
purposely runs very tight. To work around this we can reduce the RAM allocated
to the Cloud Server to give more to the hypervisor. We call this Cloud Block
Storage Optimized.`,

	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag: "UniqId",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		apiArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
		}

		var details apiTypes.CloudServerIsBlockStorageOptimized
		if err := lwCliInst.CallLwApiInto("bleed/storm/server/issbsoptimized",
			apiArgs, &details); err != nil {
			lwCliInst.Die(err)
		}

		if details.IsOptimized {
			fmt.Println("Cloud Server is Cloud Block Storage Optimized")
		} else {
			fmt.Println("Cloud Server is not Cloud Block Storage Optimized")
		}
	},
}

func init() {
	cloudServerBlockStorageOptimizedCmd.AddCommand(cloudServerBlockStorageOptimizedCheckCmd)
	cloudServerBlockStorageOptimizedCheckCmd.Flags().String("uniq-id", "", "uniq-id of Cloud Server")

	if err := cloudServerBlockStorageOptimizedCheckCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
