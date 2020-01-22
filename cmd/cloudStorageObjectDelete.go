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

var cloudStorageObjectDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an Object Store",
	Long:  `Delete an Object Store`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")

		validateFields := map[string]interface{}{
			"UniqId": uniqIdFlag,
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(fmt.Errorf("passed uniq_id [%s] failed validation: %s", uniqIdFlag, err))
		}

		apiArgs := map[string]interface{}{"uniq_id": uniqIdFlag}

		var details apiTypes.CloudObjectStoreDelete
		err := lwCliInst.CallLwApiInto("bleed/storage/objectstore/delete", apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("deleted object store %s\n", details.Deleted)
	},
}

func init() {
	cloudStorageObjectCmd.AddCommand(cloudStorageObjectDeleteCmd)
	cloudStorageObjectDeleteCmd.Flags().String("uniq_id", "",
		"uniq_id of object store to delete (see 'cloud inventory storage object list')")

	cloudStorageObjectDeleteCmd.MarkFlagRequired("uniq_id")
}
