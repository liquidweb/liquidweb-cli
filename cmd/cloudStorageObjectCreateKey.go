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

var cloudStorageObjectCreateKeyCmd = &cobra.Command{
	Use:   "createkey",
	Short: "Create a new key for the Object Store",
	Long:  `Create a new key for the Object Store`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")

		validateFields := map[interface{}]string{
			uniqIdFlag: "UniqId",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		apiArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
		}

		var details apiTypes.CloudObjectStoreKeyDetails
		err := lwCliInst.CallLwApiInto("bleed/storage/objectstore/createkey", apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Created Key:\n")
		fmt.Printf("\tUser: %s\n", details.User)
		fmt.Printf("\tAccess Key: %s\n", details.AccessKey)
		fmt.Printf("\tSecret Key: %s\n", details.SecretKey)
	},
}

func init() {
	cloudStorageObjectCmd.AddCommand(cloudStorageObjectCreateKeyCmd)
	cloudStorageObjectCreateKeyCmd.Flags().String("uniq_id", "", "uniq_id of Object Store")

	cloudStorageObjectCreateKeyCmd.MarkFlagRequired("uniq_id")
}
