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

var cloudImageRenameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Renames a Cloud Image",
	Long:  `Renames a Cloud Image.`,
	Run: func(cmd *cobra.Command, args []string) {
		imageIdFlag, _ := cmd.Flags().GetInt64("image_id")
		nameFlag, _ := cmd.Flags().GetString("name")

		if imageIdFlag == -1 || nameFlag == "" {
			lwCliInst.Die(fmt.Errorf("flags --image_id --name are required"))
		}

		apiArgs := map[string]interface{}{"id": imageIdFlag, "name": nameFlag}

		var details apiTypes.CloudImageDetails
		err := lwCliInst.CallLwApiInto("bleed/storm/image/update", apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(details)
	},
}

func init() {
	cloudImageCmd.AddCommand(cloudImageRenameCmd)

	cloudImageRenameCmd.Flags().Int64("image_id", -1,
		"id number of the image (see 'cloud inventory image list')")
	cloudImageRenameCmd.Flags().String("name", "", "new name for the Cloud Image")
}
