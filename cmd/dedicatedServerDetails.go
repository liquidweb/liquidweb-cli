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

var dedicatedServerDetailsCmdUniqIdFlag []string

var dedicatedServerDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get details of a dedicated server",
	Long:  `Get details of a dedicated server`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")

		for _, uniqId := range dedicatedServerDetailsCmdUniqIdFlag {
			validateFields := map[interface{}]interface{}{
				uniqId: "UniqId",
			}
			if err := validate.Validate(validateFields); err != nil {
				fmt.Printf("%s ... skipping\n", err)
				continue
			}

			var details apiTypes.Subaccnt
			apiArgs := map[string]interface{}{
				"uniq_id":  uniqId,
				"alsowith": []string{"categories"},
			}

			if err := lwCliInst.CallLwApiInto("bleed/asset/details", apiArgs, &details); err != nil {
				lwCliInst.Die(err)
			}

			var found bool
			for _, category := range details.Categories {
				if category == "StrictDedicated" {
					found = true
					break
				}
			}

			if !found {
				lwCliInst.Die(fmt.Errorf("UniqId [%s] is not a dedicated server", uniqId))
			}

			if jsonFlag {
				pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(details)
				if err != nil {
					panic(err)
				}
				fmt.Print(pretty)
			} else {
				fmt.Print(details)
			}
		}
	},
}

func init() {
	dedicatedServerCmd.AddCommand(dedicatedServerDetailsCmd)

	dedicatedServerDetailsCmd.Flags().Bool("json", false, "output in json format")
	dedicatedServerDetailsCmd.Flags().StringSliceVar(&dedicatedServerDetailsCmdUniqIdFlag, "uniq-id", []string{},
		"uniq-id of the dedicated server. For multiple, must be ',' separated")

	if err := dedicatedServerDetailsCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
