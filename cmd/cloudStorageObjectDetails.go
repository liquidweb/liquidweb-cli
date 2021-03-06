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

var cloudStorageObjectDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get details of a Object Store",
	Long:  `Get details of a Object Store`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag: "UniqId",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		var details apiTypes.CloudObjectStoreDetails
		apiArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
		}
		if err := lwCliInst.CallLwApiInto("bleed/storage/objectstore/details", apiArgs, &details); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(details)
	},
}

func init() {
	cloudStorageObjectCmd.AddCommand(cloudStorageObjectDetailsCmd)
	cloudStorageObjectDetailsCmd.Flags().String("uniq-id", "", "uniq-id of the object store")

	if err := cloudStorageObjectDetailsCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
