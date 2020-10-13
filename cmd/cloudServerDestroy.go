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
	"os"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/utils"
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
		forceFlag, _ := cmd.Flags().GetBool("force")

		destroyTargets := map[string]interface{}{}

		if forceFlag {
			for _, uniqId := range cloudServerDestroyCmdUniqIdFlag {
				destroyTargets[uniqId] = "" // didnt do lookup, so dont know hostname
			}
		} else {
			methodArgs := instance.AllPaginatedResultsArgs{
				Method:         "bleed/storm/server/list",
				ResultsPerPage: 100,
			}
			results, err := lwCliInst.AllPaginatedResults(&methodArgs)
			if err != nil {
				lwCliInst.Die(err)
			}

			for _, item := range results.Items {
				uniqId := cast.ToString(item["uniq_id"])
				var shouldDestroy bool
				for _, candidateUniqId := range cloudServerDestroyCmdUniqIdFlag {
					if uniqId == candidateUniqId {
						shouldDestroy = true
						break
					}
				}

				if shouldDestroy {
					destroyTargets[uniqId] = cast.ToString(item["domain"])
				}
			}

			if len(destroyTargets) == 0 {
				lwCliInst.Die(fmt.Errorf("no Cloud Servers found to destroy"))
			}

			utils.PrintYellow("DANGER! ")
			utils.PrintRed("This will destroy ALL Cloud Servers listed below:\n\n")
			for uniqId, hostname := range destroyTargets {
				fmt.Printf("\tuniq_id: %s hostname: %s\n", uniqId, hostname)
			}
			fmt.Println("")

			if proceed := dialogDesctructiveConfirmProceed(); !proceed {
				os.Exit(0)
			}
		}

		for uniqId, hostname := range destroyTargets {
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

			if hostname == "" {
				fmt.Printf("destroyed: %s\n", destroyed.Destroyed)
			} else {
				fmt.Printf("destroyed: %s (%s)\n", destroyed.Destroyed, hostname)
			}
		}
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerDestroyCmd)

	cloudServerDestroyCmd.Flags().StringSliceVar(&cloudServerDestroyCmdUniqIdFlag, "uniq-id",
		[]string{}, "uniq-ids separated by ',' of server(s) to destroy")
	cloudServerDestroyCmd.Flags().String("comment", "initiated from liquidweb-cli",
		"comment related to the cancellation")
	cloudServerDestroyCmd.Flags().String("reason", "",
		"reason for the cancellation (optional)")
	cloudServerDestroyCmd.Flags().Bool("force", false, "bypass dialog confirmation")

	if err := cloudServerDestroyCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
