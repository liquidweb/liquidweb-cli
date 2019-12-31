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

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/types/api"
)

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

		flagUniqId, _ := cmd.Flags().GetString("uniq_id")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		if flagUniqId == "" {
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
				_printCloudServerStatus(item["uniq_id"].(string), jsonOutput)
			}
		} else {
			_printCloudServerStatus(flagUniqId, jsonOutput)
		}
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerStatusCmd)

	cloudServerStatusCmd.Flags().Bool("json", false, "output in json format")
	cloudServerStatusCmd.Flags().String("uniq_id", "", "only get the status of this uniq_id")
}

func _printCloudServerStatus(uniqId string, jsonOutput bool) {
	var status apiTypes.CloudServerStatus
	if err := lwCliInst.CallLwApiInto("bleed/storm/server/status", map[string]interface{}{"uniq_id": uniqId}, &status); err != nil {
		lwCliInst.Die(err)
	}
	if jsonOutput {
		pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(status)
		if err != nil {
			lwCliInst.Die(err)
		}
		fmt.Printf(pretty)
	} else {
		fmt.Printf("uniq_id: %s\n", uniqId)
		fmt.Printf("\tstatus: %s\n", status.Status)
		if len(status.Running) > 0 {
			fmt.Printf("\tdetailed status: %s\n", status.DetailedStatus)
			fmt.Printf("\trunning: %+v\n", status.Running)
			fmt.Printf("\tprogress: %+v\n", status.Progress)
		}
	}
}
