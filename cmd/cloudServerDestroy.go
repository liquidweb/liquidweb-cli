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

var cloudServerDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a Cloud Server",
	Long: `Destroy a Cloud Server.

Kills a server. It will refund for any remaining time that has been prepaid, charge any outstanding bandwidth
charges, and then start the workflow to tear down the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		jsonFlag, _ := cmd.Flags().GetBool("json")
		commentFlag, _ := cmd.Flags().GetString("comment")
		reasonFlag, _ := cmd.Flags().GetString("reason")

		if uniqIdFlag == "" {
			lwCliInst.Die(fmt.Errorf("--uniq-id is a required flag"))
		}

		destroyArgs := map[string]interface{}{
			"uniq_id":              uniqIdFlag,
			"cancellation_comment": commentFlag,
		}

		if reasonFlag != "" {
			destroyArgs["cancellation_reason"] = reasonFlag
		}

		var destroyed apiTypes.CloudServerDestroyResponse
		if err := lwCliInst.CallLwApiInto("bleed/server/destroy", destroyArgs, &destroyed); err != nil {
			lwCliInst.Die(err)
		}

		if jsonFlag {
			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(destroyed)
			if err != nil {
				lwCliInst.Die(err)
			}
			fmt.Printf(pretty)
		} else {
			fmt.Printf("destroyed: %s\n", destroyed.Destroyed)
		}

	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerDestroyCmd)
	cloudServerDestroyCmd.Flags().Bool("json", false, "output in json format")
	cloudServerDestroyCmd.Flags().String("uniq_id", "", "uniq_id of server to destroy")
	cloudServerDestroyCmd.Flags().String("comment", "initiated from liquidweb-cli", "comment related to the cancellation")
	cloudServerDestroyCmd.Flags().String("reason", "", "reason for the cancellation (optional)")
}
