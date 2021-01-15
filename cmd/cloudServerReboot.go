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
)

var cloudServerRebootCmdUniqIdFlag []string

var cloudServerRebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot a Cloud Server",
	Long: `Reboots the cloud server as specified by the uniq-id flag.

To perform a forced a reboot, you must use --force`,
	Run: func(cmd *cobra.Command, args []string) {
		params := &instance.CloudServerRebootParams{}
		params.Force, _ = cmd.Flags().GetBool("force")

		for _, uniqId := range cloudServerRebootCmdUniqIdFlag {
			params.UniqId = uniqId

			status, err := lwCliInst.CloudServerReboot(params)
			if err != nil {
				lwCliInst.Die(err)
			}

			fmt.Print(status)
		}
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerRebootCmd)

	cloudServerRebootCmd.Flags().Bool("force", false, "perform a forced reboot")
	cloudServerRebootCmd.Flags().StringSliceVar(&cloudServerRebootCmdUniqIdFlag, "uniq-id", []string{},
		"uniq-id(s) to get status of. For multiple, must be ',' separated")

	if err := cloudServerRebootCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
