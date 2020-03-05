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

var cloudServerUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update details of your Cloud Server",
	Long: `Update details of your Cloud Server.

Update details about your server, including the backup and bandwidth plans, and the hostname
('domain') we have on file. Updating the 'domain' field will not change the actual hostname
on the server. It merely updates what our records show.

bandwidth_plan is the bandwidth plan you wish to use.  A quota of 0 indicates that you want
as-you-go, usage-based bandwidth charges.

A quota backup plan allows you to save daily backups up to your set quota.
A daily backup plan will save backups of a server up to your set days.

Either backup plan has a maximum retention of 90 days.
`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")
		hostnameFlag, _ := cmd.Flags().GetString("hostname")
		disableBackupsFlag, _ := cmd.Flags().GetBool("disable-backups")
		bandwidthQuotaFlag, _ := cmd.Flags().GetInt64("bandwidth-quota")
		backupDaysFlag, _ := cmd.Flags().GetInt64("backup-days")
		backupQuotaFlag, _ := cmd.Flags().GetInt64("backup-quota")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag: "UniqId",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		if backupDaysFlag != -1 && backupQuotaFlag != -1 {
			lwCliInst.Die(fmt.Errorf("--backup-days and --backup-quota are conflicting flags"))
		}

		if backupDaysFlag == -1 && backupQuotaFlag == -1 && !disableBackupsFlag &&
			hostnameFlag == "" && bandwidthQuotaFlag == -1 {
			lwCliInst.Die(fmt.Errorf(
				"must pass a valid flag; check 'help cloud server update' for usage"))
		}
		if disableBackupsFlag {
			if backupDaysFlag != -1 || backupQuotaFlag != -1 {
				lwCliInst.Die(fmt.Errorf("cant both enable and disable backups"))
			}
		}

		apiArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
		}

		if hostnameFlag != "" {
			apiArgs["domain"] = hostnameFlag
		}
		if bandwidthQuotaFlag != -1 {
			apiArgs["bandwidth_quota"] = bandwidthQuotaFlag
		}

		if backupDaysFlag != -1 {
			apiArgs["backup_plan"] = "daily"
			apiArgs["backup_quota"] = backupDaysFlag
		} else if backupQuotaFlag != -1 {
			apiArgs["backup_plan"] = "quota"
			apiArgs["backup_quota"] = backupQuotaFlag
		} else if disableBackupsFlag {
			apiArgs["backup_plan"] = "None"
		}

		var details apiTypes.CloudServerDetails
		if err := lwCliInst.CallLwApiInto("bleed/storm/server/update", apiArgs, &details); err != nil {
			lwCliInst.Die(err)
		}

		_printExtendedCloudServerDetails(&details)
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerUpdateCmd)

	cloudServerUpdateCmd.Flags().String("uniq-id", "", "uniq-id of the Cloud Server")
	cloudServerUpdateCmd.Flags().String("hostname", "", "hostname to set (will only update record at LiquidWeb)")
	cloudServerUpdateCmd.Flags().Int64("backup-days", -1, "Enable daily backup plan. This is the amount of days to keep a backup")
	cloudServerUpdateCmd.Flags().Int64("backup-quota", -1, "Enable quota backup plan. This is the total amount of GB to keep.")
	cloudServerUpdateCmd.Flags().Bool("disable-backups", false, "disable backups")
	cloudServerUpdateCmd.Flags().Int64("bandwidth-quota", -1,
		"bandwidth quota (0 indicates as-you-go, usage-based bandwidth charges)")

	cloudServerUpdateCmd.MarkFlagRequired("uniq-id")
}
