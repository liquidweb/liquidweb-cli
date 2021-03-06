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

var cloudImageCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Cloud Image",
	Long:  `Create a Cloud Image.`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")
		nameFlag, _ := cmd.Flags().GetString("name")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag: "UniqId",
			nameFlag:   "NonEmptyString",
		}
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		apiArgs := map[string]interface{}{"name": nameFlag, "uniq_id": uniqIdFlag}

		var details apiTypes.CloudImageCreateResponse
		err := lwCliInst.CallLwApiInto("bleed/storm/image/create", apiArgs, &details)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Creating image! %+v\n", details)
		fmt.Printf("\tthe Cloud Image will not appear in 'cloud image list' until complete\n")
	},
}

func init() {
	cloudImageCmd.AddCommand(cloudImageCreateCmd)

	cloudImageCreateCmd.Flags().String("uniq-id", "", "uniq-id of Cloud Server")
	cloudImageCreateCmd.Flags().String("name", "", "name for the Cloud Image")

	if err := cloudImageCreateCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
	if err := cloudImageCreateCmd.MarkFlagRequired("name"); err != nil {
		lwCliInst.Die(err)
	}
}
