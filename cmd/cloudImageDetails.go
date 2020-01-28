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

var cloudImageDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get details of a Cloud Image",
	Long:  `Get details of a Cloud Image`,
	Run: func(cmd *cobra.Command, args []string) {
		imageIdFlag, _ := cmd.Flags().GetInt64("image_id")

		apiArgs := map[string]interface{}{"id": imageIdFlag}

		var details apiTypes.CloudImageDetails
		err := lwCliInst.CallLwApiInto("bleed/storm/image/details", apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(details)
	},
}

func init() {
	cloudImageCmd.AddCommand(cloudImageDetailsCmd)

	cloudImageDetailsCmd.Flags().Int64("image_id", -1,
		"id number of the image (see 'cloud image list')")
	cloudImageDetailsCmd.MarkFlagRequired("image_id")
}
