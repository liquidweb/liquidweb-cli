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
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		rebuildFsFlag, _ := cmd.Flags().GetBool("rebuild-fs")
		backupIdFlag, _ := cmd.Flags().GetInt64("backup_id")

		validateFields := map[string]interface{}{
			"UniqId":        uniqIdFlag,
			"PositiveInt64": backupIdFlag,
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(fmt.Errorf("flag validation failure: %s", err))
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
		fmt.Printf("\tcheck progress with 'cloud server status --uniq_id %s'\n", uniqIdFlag)
	},
}

func init() {
	cloudBackupCmd.AddCommand(cloudBackupRestoreCmd)

	cloudBackupRestoreCmd.Flags().String("uniq_id", "", "uniq_id of Cloud Server")
	cloudBackupRestoreCmd.Flags().Int64("backup_id", -1, "id of the Cloud Backup")
	cloudBackupRestoreCmd.Flags().Bool("rebuild-fs", false, "rebuild filesystem before restoring")

	cloudBackupRestoreCmd.MarkFlagRequired("uniq_id")
	cloudBackupRestoreCmd.MarkFlagRequired("backup_id")
}
