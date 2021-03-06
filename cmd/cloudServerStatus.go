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
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/utils"
	"github.com/liquidweb/liquidweb-cli/validate"
)

var cloudServerStatusCmdUniqIdFlag []string

var cloudServerStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get status of Cloud Server(s)",
	Long: `Get status of Cloud Server(s).

Gets the current status of a server, in relation to what processes are being run on its behalf.

The list of available statuses at this time are:

	Building
	Cloning
	Resizing
	Moving
	Booting
	Stopping
	Restarting
	Rebooting
	Shutting Down
	Restoring Backup
	Creating Image
	Deleting Image
	Restoring Image
	Re-Imaging
	Updating Firewall
	Updating Network
	Adding IPs
	Removing IP
	Destroying

If nothing is currently running, only the 'status' field will be returned with one of the following statuses:

	Failed
	Provisioning
	Running
	Shutdown
	Stopped
`,
	Run: func(cmd *cobra.Command, args []string) {
		watchFlag, _ := cmd.Flags().GetBool("watch")
		intervalFlag, _ := cmd.Flags().GetInt("interval")

		if watchFlag {
			if intervalFlag <= 0 {
				lwCliInst.Die(fmt.Errorf("You must specify an interval greater than zero."))
			}

			for {
				fmt.Println("\nDisplaying server status (CTRL-C to exit):")
				displayCloudSErverStatus(cloudServerStatusCmdUniqIdFlag)
				time.Sleep(time.Duration(intervalFlag) * time.Second)
			}
		} else {
			displayCloudSErverStatus(cloudServerStatusCmdUniqIdFlag)
		}
	},
}

func displayCloudSErverStatus(uniqIdList []string) {

	if len(uniqIdList) == 0 {
		// fetch status of all cloud servers on account
		methodArgs := instance.AllPaginatedResultsArgs{
			Method:         "bleed/storm/server/list",
			ResultsPerPage: 100,
		}
		results, err := lwCliInst.AllPaginatedResults(&methodArgs)
		if err != nil {
			lwCliInst.Die(err)
		}

		for _, item := range results.Items {
			var details apiTypes.CloudServerDetails
			if err := instance.CastFieldTypes(item, &details); err != nil {
				lwCliInst.Die(err)
			}

			_printCloudServerStatus(details.UniqId, details.Domain)
		}
	} else {
		for _, uid := range uniqIdList {
			validateFields := map[interface{}]interface{}{
				uid: "UniqId",
			}
			if err := validate.Validate(validateFields); err != nil {
				fmt.Printf("%s ... skipping\n", err)
				continue
			}
			_printCloudServerStatus(uid, "")
		}
	}

}

func init() {
	cloudServerCmd.AddCommand(cloudServerStatusCmd)

	cloudServerStatusCmd.Flags().StringSliceVar(&cloudServerStatusCmdUniqIdFlag, "uniq-id", []string{},
		"uniq-id(s) to get status of. For multiple, must be ',' separated")
	cloudServerStatusCmd.Flags().Bool("watch", false, "continue to redisplay status at --interval")
	cloudServerStatusCmd.Flags().Int("interval", 10, "the interval (in seconds) to fetch the status when --watch is specified")
}

func _printCloudServerStatus(uniqId string, domain string) {
	var status apiTypes.CloudServerStatus
	if err := lwCliInst.CallLwApiInto("bleed/storm/server/status", map[string]interface{}{"uniq_id": uniqId},
		&status); err != nil {
		lwCliInst.Die(err)
	}

	if domain == "" {
		var details apiTypes.CloudServerDetails
		if err := lwCliInst.CallLwApiInto("bleed/storm/server/details",
			map[string]interface{}{"uniq_id": uniqId}, &details); err != nil {
			lwCliInst.Die(err)
		}
		domain = details.Domain
	}

	utils.PrintTeal("UniqId: %s\n", uniqId)
	fmt.Printf("\tdomain: %s\n", domain)
	if strings.ToUpper(status.Status) == "RUNNING" {
		fmt.Printf("\tstatus: %s\n", status.Status)
	} else {
		utils.PrintYellow("\tstatus: %s\n", status.Status)
	}
	if len(status.Running) > 0 {
		fmt.Printf("\tdetailed status: %s\n", status.DetailedStatus)
		fmt.Printf("\trunning: %+v\n", status.Running)
		fmt.Printf("\tprogress: %+v\n", status.Progress)
	}
}
