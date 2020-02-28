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

var cloudPrivateParentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Private Parent",
	Long: `Create a new Private Parent.

A Private Parent is a physical hypervisor node that you fully own. No one else but you
will be able to provision Cloud Servers on a Private Parent. In addition, with Private
Parents you have total control of how many instances can live on the Private Parent,
as well as how many resources each Cloud Server gets.

Private Parents must use a config of category 'bare-metal' or 'bare-metal-r'. For a list
of configs, check 'cloud server options --configs'.`,
	Run: func(cmd *cobra.Command, args []string) {
		nameFlag, _ := cmd.Flags().GetString("name")
		configIdFlag, _ := cmd.Flags().GetInt64("config-id")
		zoneFlag, _ := cmd.Flags().GetInt64("zone")

		validateFields := map[interface{}]interface{}{
			zoneFlag:     "PositiveInt64",
			configIdFlag: "PositiveInt64",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		apiArgs := map[string]interface{}{
			"domain":    nameFlag,
			"config_id": configIdFlag,
			"zone":      zoneFlag,
		}

		var details apiTypes.CloudPrivateParentDetails
		if err := lwCliInst.CallLwApiInto("bleed/storm/private/parent/create",
			apiArgs, &details); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Private Parent with name [%s] uniq-id [%s] created!\n", details.Domain, details.UniqId)
		fmt.Printf("\tYou can now provision Cloud Servers on this Private Parent. See 'help cloud server create'\n")
	},
}

func init() {
	cloudPrivateParentCmd.AddCommand(cloudPrivateParentCreateCmd)

	cloudPrivateParentCreateCmd.Flags().Int64("config-id", -1, "config-id (category must be bare-metal or bare-metal-r)")
	cloudPrivateParentCreateCmd.Flags().String("name", "", "name for your Private Parent")
	cloudPrivateParentCreateCmd.Flags().Int64("zone", -1, "id number of the zone to provision the Private Parent in ('cloud server options --zones')")

	cloudPrivateParentCreateCmd.MarkFlagRequired("config-id")
	cloudPrivateParentCreateCmd.MarkFlagRequired("zone")
	cloudPrivateParentCreateCmd.MarkFlagRequired("name")
}
