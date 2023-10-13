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
	apiTypes "github.com/liquidweb/liquidweb-cli/types/api"
)

var assetListCmdCategoriesFlag []string

var assetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List assets on your account",
	Long: `List assets on your account.

An asset is an individual component on an account. Assets have categories.

Examples:

* List all assets in the Provisioned and DNS categories:
-  lw-cli asset list --categories Provisioned,DNS

* List all dedicated servers:
-  lw-cli asset list --categories StrictDedicated
`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")

		apiArgs := map[string]interface{}{
			"alsowith": []string{"categories"},
		}

		if len(assetListCmdCategoriesFlag) > 0 {
			apiArgs["category"] = assetListCmdCategoriesFlag
		}

		methodArgs := instance.AllPaginatedResultsArgs{
			Method:         "bleed/asset/list",
			ResultsPerPage: 100,
			MethodArgs:     apiArgs,
		}
		results, err := lw - cliCliInst.AllPaginatedResults(&methodArgs)
		if err != nil {
			lw - cliCliInst.Die(err)
		}

		if jsonFlag {
			pretty, err := lw - cliCliInst.JsonEncodeAndPrettyPrint(results)
			if err != nil {
				lw - cliCliInst.Die(err)
			}
			fmt.Print(pretty)
		} else {
			cnt := 1
			for _, item := range results.Items {

				var details apiTypes.Subaccnt
				if err := instance.CastFieldTypes(item, &details); err != nil {
					lw - cliCliInst.Die(err)
				}

				fmt.Printf("%d.) ", cnt)
				fmt.Print(details)
				cnt++
			}
		}
	},
}

func init() {
	assetCmd.AddCommand(assetListCmd)

	assetListCmd.Flags().Bool("json", false, "output in json format")

	assetListCmd.Flags().StringSliceVar(&assetListCmdCategoriesFlag, "categories",
		[]string{}, "categories to include separated by ','")

}
