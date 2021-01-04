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
	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/validate"
)

var cloudNetworkPrivateDetailsCmdUniqIdFlag []string

var cloudNetworkPrivateDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get Private Network details for a single or all Cloud Server(s)",
	Long: `Get Private Network details for a single or all Cloud Server(s)

Private networking provides the option for several Cloud Servers to contact each other
via a network interface that is:

A) not publicly routable
B) not metered for bandwidth.

Applications that communicate internally will frequently use this for both security
and cost-savings.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var uniqIds []string

		if len(cloudNetworkPrivateDetailsCmdUniqIdFlag) == 0 {
			methodArgs := instance.AllPaginatedResultsArgs{
				Method:         "bleed/storm/server/list",
				ResultsPerPage: 100,
			}
			results, err := lwCliInst.AllPaginatedResults(&methodArgs)
			if err != nil {
				lwCliInst.Die(err)
			}
			for _, item := range results.Items {
				var cs apiTypes.CloudServerDetails
				if err := instance.CastFieldTypes(item, &cs); err != nil {
					lwCliInst.Die(err)
				}
				uniqIds = append(uniqIds, cs.UniqId)
			}
		} else {
			uniqIds = cloudNetworkPrivateDetailsCmdUniqIdFlag
		}

		for _, uniqId := range uniqIds {
			validateFields := map[interface{}]interface{}{
				uniqId: "UniqId",
			}
			if err := validate.Validate(validateFields); err != nil {
				lwCliInst.Die(err)
			}

			apiArgs := map[string]interface{}{"uniq_id": uniqId}

			var details apiTypes.CloudNetworkPrivateGetIpResponse
			if err := lwCliInst.CallLwApiInto("bleed/network/private/getip", apiArgs, &details); err != nil {
				lwCliInst.Die(err)
			}

			fmt.Print(details)
		}
	},
}

func init() {
	cloudNetworkPrivateCmd.AddCommand(cloudNetworkPrivateDetailsCmd)
	cloudNetworkPrivateDetailsCmd.Flags().StringSliceVar(&cloudNetworkPrivateDetailsCmdUniqIdFlag, "uniq-id",
		[]string{}, "uniq-ids separated by ',' of Cloud Servers to fetch private networking details for")
}
