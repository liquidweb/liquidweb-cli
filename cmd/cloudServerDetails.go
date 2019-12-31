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
)

var cloudServerDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get details of a Cloud Server",
	Long: `Get details of a Cloud Server.

You can check this methods API documentation for what the returned fields mean:

https://cart.liquidweb.com/storm/api/docs/bleed/Storm/Server.html#method_details
`,
	Run: func(cmd *cobra.Command, args []string) {

		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		if uniqIdFlag == "" {
			lwCliInst.Die(fmt.Errorf("--uniq_id is a required flag"))
		}

		details, err := lwCliInst.LwApiClient.Call("bleed/storm/server/details", map[string]interface{}{"uniq_id": uniqIdFlag})
		if err != nil {
			lwCliInst.Die(err)
		}

		if jsonFlag {
			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(details)
			if err != nil {
				lwCliInst.Die(err)
			}
			fmt.Printf(pretty)
			os.Exit(0)
		}

		for key, value := range details.(map[string]interface{}) {
			if key == "accnt" {
				value = cast.ToInt64(value)
			}
			fmt.Printf("%s: %+v\n", key, value)
		}
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerDetailsCmd)

	cloudServerDetailsCmd.Flags().Bool("json", false, "output in json format")
	cloudServerDetailsCmd.Flags().String("uniq_id", "", "get details of this uniq_id")
}
