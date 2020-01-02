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
)

var cloudServerShutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Shutdown a Cloud Server",
	Long: `Shutdown a Cloud Server.

Stop a server. The 'force' flag will do a hard stop of the server from the parent server. Otherwise, it
will issue a halt command to the server and shutdown normally.`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		if uniqIdFlag == "" {
			lwCliInst.Die(fmt.Errorf("--uniq-id is a required flag"))
		}

		shutdownArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
		}

		result, err := lwCliInst.LwApiClient.Call("bleed/server/shutdown", shutdownArgs)
		if err != nil {
			lwCliInst.Die(err)
		}

		if jsonFlag {
			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(result)
			if err != nil {
				lwCliInst.Die(err)
			}
			fmt.Printf(pretty)
		} else {
			fmt.Printf("shutdown: %s\n", result.(map[string]interface{})["shutdown"])
		}

	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerShutdownCmd)
	cloudServerShutdownCmd.Flags().Bool("json", false, "output in json format")
	cloudServerShutdownCmd.Flags().String("uniq_id", "", "uniq_id of server to shutdown")
}