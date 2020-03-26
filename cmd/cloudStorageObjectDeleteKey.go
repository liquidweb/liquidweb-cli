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
	"github.com/liquidweb/liquidweb-cli/validate"
)

var cloudStorageObjectDeleteKeyCmd = &cobra.Command{
	Use:   "deletekey",
	Short: "Delete a key from the Object Store",
	Long:  `Delete a key from the Object Store`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")
		accessKeyFlag, _ := cmd.Flags().GetString("access-key")
		forceFlag, _ := cmd.Flags().GetBool("force")

		// if force flag wasn't passed
		if !forceFlag {
			// exit if user didn't consent
			if proceed := dialogDesctructiveConfirmProceed(); !proceed {
				os.Exit(0)
			}
		}

		validateFields := map[interface{}]interface{}{
			uniqIdFlag:    "UniqId",
			accessKeyFlag: "NonEmptyString",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
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
	cloudStorageObjectDeleteKeyCmd.Flags().String("uniq-id", "", "uniq-id of Object Store")
	cloudStorageObjectDeleteKeyCmd.Flags().String("access-key", "", "the access key to remove from the Object Store")
	cloudStorageObjectDeleteKeyCmd.Flags().Bool("force", false, "bypass dialog confirmation")

	if err := cloudStorageObjectDeleteKeyCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
	if err := cloudStorageObjectDeleteKeyCmd.MarkFlagRequired("access-key"); err != nil {
		lwCliInst.Die(err)
	}
}
