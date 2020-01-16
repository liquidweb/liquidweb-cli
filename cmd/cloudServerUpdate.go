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

var cloudServerUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update details of your Cloud Server",
	Long: `Update details of your Cloud Server.

Update details about your server, including the backup and bandwidth plans, and the hostname
('domain') we have on file. Updating the 'domain' field will not change the actual hostname
on the server. It merely updates what our records show.

bandwidth_plan is the bandwidth plan you wish to use.  A quota of 0 indicates that you want
as-you-go, usage-based bandwidth charges.`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		hostnameFlag, _ := cmd.Flags().GetString("hostname")
		backupPlanFlag, _ := cmd.Flags().GetString("backup-plan")
		disableBackupsFlag, _ := cmd.Flags().GetBool("disable-backups")
		bandwidthQuotaFlag, _ := cmd.Flags().GetInt64("bandwidth-quota")
		backupQuotaFlag, _ := cmd.Flags().GetInt64("backup-quota")

		if uniqIdFlag == "" {
			lwCliInst.Die(fmt.Errorf("flag --uniq_id is required"))
		}

		if backupPlanFlag == "Quota" {
			if backupQuotaFlag == -1 {
				lwCliInst.Die(fmt.Errorf("cannot enable Quota backups without --backup-quota"))
			}
		}

		if backupPlanFlag == "" && !disableBackupsFlag && hostnameFlag == "" &&
			bandwidthQuotaFlag == -1 {
			lwCliInst.Die(fmt.Errorf(
				"must pass one of: enable-backups disable-backups hostname bandwidth-quota"))
		}
		if backupPlanFlag != "" && disableBackupsFlag {
			lwCliInst.Die(fmt.Errorf("cant both enable and disable backups"))
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

		if backupPlanFlag != "" {
			apiArgs["backup_plan"] = backupPlanFlag
			if backupPlanFlag == "Quota" {
				apiArgs["backup_quota"] = backupQuotaFlag
			}
		}

		if disableBackupsFlag {
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

	cloudServerUpdateCmd.Flags().String("uniq_id", "", "uniq_id of the Cloud Server")
	cloudServerUpdateCmd.Flags().String("hostname", "", "hostname to set")
	cloudServerUpdateCmd.Flags().String("backup-plan", "", "Name of the backup plan")
	cloudServerUpdateCmd.Flags().Int64("backup-quota", -1, "Quota to set for Quota type backup-plan")
	cloudServerUpdateCmd.Flags().Bool("disable-backups", false, "disable backups")
	cloudServerUpdateCmd.Flags().Int64("bandwidth-quota", -1,
		"bandwidth quota (0 indicates as-you-go, usage-based bandwidth charges)")
}
