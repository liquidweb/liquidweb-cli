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

var cloudServerDestroyCmdUniqIdFlag []string

var cloudServerDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a Cloud Server",
	Long: `Destroy a Cloud Server.

Kills a server. It will refund for any remaining time that has been prepaid, charge
any outstanding bandwidth charges, and then start the workflow to tear down the
server.`,
	Run: func(cmd *cobra.Command, args []string) {
		commentFlag, _ := cmd.Flags().GetString("comment")
		reasonFlag, _ := cmd.Flags().GetString("reason")

		for _, uniqId := range cloudServerDestroyCmdUniqIdFlag {

			validateFields := map[interface{}]interface{}{
				uniqId: "UniqId",
			}

			if err := validate.Validate(validateFields); err != nil {
				fmt.Printf("%s ... skipping\n", err)
				continue
			}

			destroyArgs := map[string]interface{}{
				"uniq_id":              uniqId,
				"cancellation_comment": commentFlag,
			}

			if reasonFlag != "" {
				destroyArgs["cancellation_reason"] = reasonFlag
			}

			var destroyed apiTypes.CloudServerDestroyResponse
			err := lwCliInst.CallLwApiInto("bleed/server/destroy", destroyArgs, &destroyed)
			if err != nil {
				lwCliInst.Die(err)
			}

			fmt.Printf("destroyed: %s\n", destroyed.Destroyed)
		}
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerDestroyCmd)

	cloudServerDestroyCmd.Flags().StringSliceVar(&cloudServerDestroyCmdUniqIdFlag, "uniq_id",
		[]string{}, "uniq_ids separated by ',' of server(s) to destroy")
	cloudServerDestroyCmd.Flags().String("comment", "initiated from liquidweb-cli",
		"comment related to the cancellation")
	cloudServerDestroyCmd.Flags().String("reason", "",
		"reason for the cancellation (optional)")

	cloudServerDestroyCmd.MarkFlagRequired("uniq_id")
}
