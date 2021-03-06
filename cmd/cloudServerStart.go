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

var cloudServerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Cloud Server",
	Long: `Start a Cloud Server.

Boot a server. If the server is already running, this will do nothing.`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag: "UniqId",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		startArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
		}

		var result apiTypes.CloudServerStartResponse
		if err := lwCliInst.CallLwApiInto("bleed/server/start", startArgs, &result); err != nil {
			lwCliInst.Die(err)
		}

		if jsonFlag {
			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(result)
			if err != nil {
				lwCliInst.Die(err)
			}
			fmt.Printf(pretty)
		} else {
			fmt.Printf("started: %s\n", result.Started)
		}

	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerStartCmd)
	cloudServerStartCmd.Flags().Bool("json", false, "output in json format")
	cloudServerStartCmd.Flags().String("uniq-id", "", "uniq-id of server to start")

	if err := cloudServerStartCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
