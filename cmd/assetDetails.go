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

var assetDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get details of a specific asset",
	Long: `Get details of a specific asset.

An asset is an individual product on an account. Assets have categories.
`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")

		var details apiTypes.Subaccnt
		apiArgs := map[string]interface{}{
			"uniq_id":  uniqIdFlag,
			"alsowith": []string{"categories"},
		}

		if err := lwCliInst.CallLwApiInto("bleed/asset/details", apiArgs, &details); err != nil {
			lwCliInst.Die(err)
		}

		if jsonFlag {
			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(details)
			if err != nil {
				lwCliInst.Die(err)
			}
			fmt.Printf(pretty)
		} else {
			fmt.Print(details)
		}
	},
}

func init() {
	assetCmd.AddCommand(assetDetailsCmd)

	assetDetailsCmd.Flags().Bool("json", false, "output in json format")
	assetDetailsCmd.Flags().String("uniq-id", "", "uniq-id of the asset")

	if err := assetDetailsCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
