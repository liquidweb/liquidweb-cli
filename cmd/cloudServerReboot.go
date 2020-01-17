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

var cloudServerRebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot a Cloud Server",
	Long: `Reboots the cloud server as specified by the uniq_id flag.

To perform a forced a reboot, you must use --force`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqId, _ := cmd.Flags().GetString("uniq_id")
		jsonOutput, _ := cmd.Flags().GetBool("json")
		force, _ := cmd.Flags().GetBool("force")

		var resp apiTypes.CloudServerRebootResponse
		if err := lwCliInst.CallLwApiInto("bleed/storm/server/reboot", map[string]interface{}{
			"uniq_id": uniqId, "force": force}, &resp); err != nil {
			lwCliInst.Die(err)
		}

		if jsonOutput {
			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(resp)
			if err != nil {
				lwCliInst.Die(err)
			}
			fmt.Printf(pretty)
		} else {
			fmt.Printf("Rebooted: %s\n", resp.Rebooted)
		}
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerRebootCmd)

	cloudServerRebootCmd.Flags().Bool("json", false, "output in json format")
	cloudServerRebootCmd.Flags().Bool("force", false, "perform a forced reboot")
	cloudServerRebootCmd.Flags().String("uniq_id", "", "uniq_id of server to reboot")

	cloudServerRebootCmd.MarkFlagRequired("uniq_id")
}
