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

var cloudServerShutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Shutdown a Cloud Server",
	Long: `Shutdown a Cloud Server.

Stop a server. The 'force' flag will do a hard stop of the server from the parent server. Otherwise, it
will issue a halt command to the server and shutdown normally.`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")
		forceFlag, _ := cmd.Flags().GetBool("force")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag: "UniqId",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		shutdownArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
			//"force":   forceFlag,
		}
		// conditionally adding to workaround bug in api method..
		// TODO delete and just always add once bug addressed
		if forceFlag {
			shutdownArgs["force"] = forceFlag
		}

		var result apiTypes.CloudServerShutdownResponse
		if err := lwCliInst.CallLwApiInto("bleed/server/shutdown", shutdownArgs, &result); err != nil {
			lwCliInst.Die(err)
		}

		if jsonFlag {
			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(result)
			if err != nil {
				lwCliInst.Die(err)
			}
			fmt.Printf(pretty)
		} else {
			fmt.Printf("shutdown: %s\n", result.Shutdown)
		}

	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerShutdownCmd)
	cloudServerShutdownCmd.Flags().Bool("json", false, "output in json format")
	cloudServerShutdownCmd.Flags().String("uniq-id", "", "uniq-id of server to shutdown")
	cloudServerShutdownCmd.Flags().Bool("force", false, "force shutdown server")

	if err := cloudServerShutdownCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
