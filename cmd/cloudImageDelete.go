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

	"github.com/liquidweb/liquidweb-cli/types/api"
)

var cloudImageDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Cloud Image",
	Long:  `Delete a Cloud Image`,
	Run: func(cmd *cobra.Command, args []string) {
		imageIdFlag, _ := cmd.Flags().GetInt64("image-id")
		forceFlag, _ := cmd.Flags().GetBool("force")

		// if force flag wasn't passed
		if !forceFlag {
			// exit if user didn't consent
			if proceed := dialogDesctructiveConfirmProceed(); !proceed {
				os.Exit(0)
			}
		}

		apiArgs := map[string]interface{}{"id": imageIdFlag}

		var details apiTypes.CloudImageDeleteResponse
		err := lwCliInst.CallLwApiInto("bleed/storm/image/delete", apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Deleted image %d\n", details.Deleted)
	},
}

func init() {
	cloudImageCmd.AddCommand(cloudImageDeleteCmd)

	cloudImageDeleteCmd.Flags().Int64("image-id", -1,
		"id number of the image (see 'cloud image list')")
	cloudImageDeleteCmd.Flags().Bool("force", false, "bypass dialog confirmation")
	if err := cloudImageDeleteCmd.MarkFlagRequired("image-id"); err != nil {
		lwCliInst.Die(err)
	}
}
