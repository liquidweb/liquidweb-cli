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

var cloudNetworkPrivateAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a Cloud Server to a Private Network",
	Long: `Attach a Cloud Server to a Private Network

Private networking provides the option for several Cloud Servers to contact each other
via a network interface that is:

A) not publicly routable
B) not metered for bandwidth.

Applications that communicate internally will frequently use this for both security
and cost-savings.
`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")

		if uniqIdFlag == "" {
			lwCliInst.Die(fmt.Errorf("flag --uniq_id is required"))
		}

		apiArgs := map[string]interface{}{"uniq_id": uniqIdFlag}

		var attachedDetails apiTypes.CloudNetworkPrivateIsAttachedResponse
		err := lwCliInst.CallLwApiInto("bleed/network/private/isattached", apiArgs, &attachedDetails)
		if err != nil {
			lwCliInst.Die(err)
		}
		if attachedDetails.IsAttached {
			lwCliInst.Die(fmt.Errorf("Cloud Server is already attached to the Private Network"))
		}

		var details apiTypes.CloudNetworkPrivateAttachResponse
		err = lwCliInst.CallLwApiInto("bleed/network/private/attach", apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Attaching %s to private network\n", details.Attached)
		fmt.Printf("\n\nYou can check progress with 'cloud server status --uniq_id %s'\n", uniqIdFlag)
	},
}

func init() {
	cloudNetworkPrivateCmd.AddCommand(cloudNetworkPrivateAttachCmd)
	cloudNetworkPrivateAttachCmd.Flags().String("uniq_id", "", "uniq_id of the Cloud Server")
}