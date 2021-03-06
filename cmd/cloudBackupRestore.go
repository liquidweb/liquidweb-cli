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

var cloudBackupRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a Cloud Backup on a Cloud Server",
	Long:  `Restore a Cloud Backup on a Cloud Server.`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")
		rebuildFsFlag, _ := cmd.Flags().GetBool("rebuild-fs")
		backupIdFlag, _ := cmd.Flags().GetInt64("backup-id")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag:   "UniqId",
			backupIdFlag: "PositiveInt64",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		apiArgs := map[string]interface{}{"id": backupIdFlag, "uniq_id": uniqIdFlag}
		if rebuildFsFlag {
			apiArgs["force"] = 1
		}

		var details apiTypes.CloudBackupRestoreResponse
		err := lwCliInst.CallLwApiInto("bleed/storm/backup/restore", apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Restoring backup! %+v\n", details)
		fmt.Printf("\tcheck progress with 'cloud server status --uniq-id %s'\n", uniqIdFlag)
	},
}

func init() {
	cloudBackupCmd.AddCommand(cloudBackupRestoreCmd)

	cloudBackupRestoreCmd.Flags().String("uniq-id", "", "uniq-id of Cloud Server")
	cloudBackupRestoreCmd.Flags().Int64("backup-id", -1, "id of the Cloud Backup")
	cloudBackupRestoreCmd.Flags().Bool("rebuild-fs", false, "rebuild filesystem before restoring")

	if err := cloudBackupRestoreCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
	if err := cloudBackupRestoreCmd.MarkFlagRequired("backup-id"); err != nil {
		lwCliInst.Die(err)
	}
}
