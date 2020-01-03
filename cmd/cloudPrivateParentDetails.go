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

var cloudPrivateParentDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get details of a Private Parent",
	Long: `Get details of a Private Parent.

A Private Parent is a physical hypervisor node that you fully own. No one else but you
will be able to provision Cloud Servers on a Private Parent. In addition, with Private
Parents you have total control of how many instances can live on the Private Parent,
as well as how many resources each Cloud Server gets.`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")

		if uniqIdFlag == "" {
			lwCliInst.Die(fmt.Errorf("flag --uniq_id is required"))
		}

		apiArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
		}

		var details apiTypes.CloudPrivateParentDetails
		if err := lwCliInst.CallLwApiInto("bleed/storm/private/parent/details",
			apiArgs, &details); err != nil {
			lwCliInst.Die(err)
		}

		_printPrivateParentDetailsFromDetailsStruct(&details)
	},
}

func init() {
	cloudPrivateParentCmd.AddCommand(cloudPrivateParentDetailsCmd)

	cloudPrivateParentDetailsCmd.Flags().String("uniq_id", "", "uniq_id of the private parent to get details of")
}
