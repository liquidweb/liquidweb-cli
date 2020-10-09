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
	"errors"
	"fmt"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/utils"
	"github.com/liquidweb/liquidweb-cli/validate"
)

var cloudServerResizeExpectationCmd = &cobra.Command{
	Use:   "resize-expectation",
	Short: "Determine if a Cloud Server can be resized without downtime",
	Long: `This command can be used to determine if a Cloud Server can be resized to the requested
config-id without downtime.

Depending on inventory and desired config-id (configuration) the resize could either
require a reboot to complete, or be performed entirely live. The intention of this
command is to provide the user with a sane expectation ahead of making the resize
request.

If there is no inventory available, an exception will be raised.

Its important to note, this command will *not* make any changes to your Cloud Server.
This command is purely for information gathering.
`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")
		privateParentFlag, _ := cmd.Flags().GetString("private-parent")
		diskFlag, _ := cmd.Flags().GetInt64("disk")
		memoryFlag, _ := cmd.Flags().GetInt64("memory")
		vcpuFlag, _ := cmd.Flags().GetInt64("vcpu")
		configIdFlag, _ := cmd.Flags().GetInt64("config-id")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag: "UniqId",
		}

		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		if privateParentFlag != "" && configIdFlag != -1 {
			lwCliInst.Die(errors.New("cant pass both --config-id and --private-parent flags"))
		}
		if privateParentFlag == "" && configIdFlag == -1 {
			lwCliInst.Die(errors.New("must pass --config-id or --private-parent"))
		}

		apiArgs := map[string]interface{}{
			"uniq_id":   uniqIdFlag,
			"config_id": configIdFlag,
		}

		// if private parent, add args
		if privateParentFlag != "" {
			if memoryFlag <= 0 && diskFlag <= 0 && vcpuFlag <= 0 {
				lwCliInst.Die(errors.New("when --private-parent , at least one of --memory --disk --vcpu are required"))
			}

			privateParentUniqId, err := lwCliInst.DerivePrivateParentUniqId(privateParentFlag)
			if err != nil {
				lwCliInst.Die(err)
			}

			var cloudServerDetails apiTypes.CloudServerDetails
			if err = lwCliInst.CallLwApiInto(
				"bleed/storm/server/details",
				map[string]interface{}{
					"uniq_id": uniqIdFlag,
				}, &cloudServerDetails); err != nil {
				lwCliInst.Die(err)
			}

			apiArgs["config_id"] = 0
			apiArgs["private_parent"] = privateParentUniqId
			apiArgs["disk"] = cloudServerDetails.DiskSpace
			apiArgs["memory"] = cloudServerDetails.Memory
			apiArgs["vcpu"] = cloudServerDetails.Vcpu

			if diskFlag > 0 {
				apiArgs["disk"] = diskFlag
			}
			if vcpuFlag > 0 {
				apiArgs["vcpu"] = vcpuFlag
			}
			if memoryFlag > 0 {
				apiArgs["memory"] = memoryFlag
			}
		}

		expectationInter, err := lwCliInst.LwCliApiClient.Call("bleed/storm/server/resizePlan", apiArgs)
		if err != nil {
			lwCliInst.Die(fmt.Errorf("ERROR: %s", err))
		}
		expectation, ok := expectationInter.(map[string]interface{})
		if !ok {
			lwCliInst.Die(errors.New("returned an unexpected structure"))
		}

		memoryDifference := cast.ToInt(expectation["memoryDifference"])
		diskDifference := cast.ToInt(expectation["diskDifference"])
		vcpuDifference := cast.ToInt(expectation["vcpuDifference"])

		utils.PrintGreen("Configuration is available\n\n")

		fmt.Print("Resource Changes: Disk [")
		if diskDifference == 0 {
			fmt.Printf("%d] ", diskDifference)
		} else if diskDifference >= 0 {
			utils.PrintGreen("%d] ", diskDifference)
		} else {
			utils.PrintRed("%d] ", diskDifference)
		}

		fmt.Print("Memory [")
		if memoryDifference == 0 {
			fmt.Printf("%d] ", memoryDifference)
		} else if memoryDifference >= 0 {
			utils.PrintGreen("%d] ", memoryDifference)
		} else {
			utils.PrintRed("%d] ", memoryDifference)
		}

		fmt.Print("Vcpu [")
		if vcpuDifference == 0 {
			fmt.Printf("%d]\n", vcpuDifference)
		} else if vcpuDifference >= 0 {
			utils.PrintGreen("%d]\n", vcpuDifference)
		} else {
			utils.PrintRed("%d]\n", vcpuDifference)
		}
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerResizeExpectationCmd)

	cloudServerResizeExpectationCmd.Flags().String("uniq-id", "", "uniq-id of Cloud Server")

	cloudServerResizeExpectationCmd.Flags().String("private-parent", "",
		"name or uniq-id of the Private Parent (see: 'cloud private-parent list')")
	cloudServerResizeExpectationCmd.Flags().Int64("disk", -1, "diskspace for the Cloud Server (when private-parent)")
	cloudServerResizeExpectationCmd.Flags().Int64("memory", -1, "memory for the Cloud Server (when private-parent)")
	cloudServerResizeExpectationCmd.Flags().Int64("vcpu", -1, "vcpus for the Cloud Server (when private-parent)")

	cloudServerResizeExpectationCmd.Flags().Int64("config-id", -1,
		"config-id to check availability for (when !private-parent) (see: 'cloud server options --configs')")

	if err := cloudServerResizeExpectationCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
