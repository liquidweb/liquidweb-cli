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

var cloudStorageObjectDeleteKeyCmd = &cobra.Command{
	Use:   "deletekey",
	Short: "Delete a key from the Object Store",
	Long:  `Delete a key from the Object Store`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		accessKeyFlag, _ := cmd.Flags().GetString("access-key")

		if uniqIdFlag == "" || accessKeyFlag == "" {
			lwCliInst.Die(fmt.Errorf("flags --uniq_id --access-key are required"))
		}

		apiArgs := map[string]interface{}{
			"uniq_id":    uniqIdFlag,
			"access_key": accessKeyFlag,
		}

		var details apiTypes.CloudObjectStoreDeleteKey
		err := lwCliInst.CallLwApiInto("bleed/storage/objectstore/deletekey", apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("deleted key %s\n", details.Deleted)
	},
}

func init() {
	cloudStorageObjectCmd.AddCommand(cloudStorageObjectDeleteKeyCmd)
	cloudStorageObjectDeleteKeyCmd.Flags().String("uniq_id", "", "uniq_id of Object Store")
	cloudStorageObjectDeleteKeyCmd.Flags().String("access-key", "", "the access key to remove from the Object Store")
}
