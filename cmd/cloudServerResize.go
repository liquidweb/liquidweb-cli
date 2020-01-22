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

var cloudServerResizeCmd = &cobra.Command{
	Use:   "resize",
	Short: "Resize a Cloud Server",
	Long: `Resize a Cloud Server.

Resize a Cloud Server to a new config. Available config_ids can be found in
'cloud server options --configs'

You will be billed for the prorated difference of the price of the new
configuration compared to the price of your old configuration. The
difference will be refunded or charged at your next billing date.

When resizing to a larger configuration, the filesystem resize operation can
be skipped by passing the 'skip-fs-resize' flag. In this case, the storage
associated with the new configuration is allocated, but not available
immediately. The filesystem resize can be scheduled with the support team,
or you can do it yourself.

Skipping the filesystem resize is only possible when moving to a larger
configuration. This option has no effect if moving to the same or smaller
configuration.

If this is a resize of a Cloud Server on a private parent, pass --private-parent
with a value of either the name of the private parent, or the private parents
uniq_id. When passing --private-parent, at least one of the following flags
are required:

  --diskspace
  --memory
  --vcpu

Downtime Expectations:

When resizing a Cloud Server on a private parent, you can add memory or vcpu(s)
without downtime. If you change the diskspace however, then a reboot will be
required.

When resizing a Cloud Server that isn't on a private parent, there will be one
reboot during the resize. The only case there will be two reboots is when
going to a config with more diskspace, and --skip-fs-resize wasn't passed.

During all resizes, the Cloud Server is online as the disk synchronizes.
`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		diskspaceFlag, _ := cmd.Flags().GetInt64("diskspace")
		configIdFlag, _ := cmd.Flags().GetInt64("config_id")
		memoryFlag, _ := cmd.Flags().GetInt64("memory")
		skipFsResizeFlag, _ := cmd.Flags().GetBool("skip-fs-resize")
		vcpuFlag, _ := cmd.Flags().GetInt64("vcpu")
		privateParentFlag, _ := cmd.Flags().GetString("private-parent")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag: "UniqId",
		}
		// must validate UniqId now because we call api methods with this uniq_id before below validate
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		// convert bool to int for api
		skipFsResizeInt := 0
		if skipFsResizeFlag {
			skipFsResizeInt = 1
		}

		if configIdFlag == -1 && privateParentFlag == "" {
			lwCliInst.Die(fmt.Errorf("flag --config_id required when --private-parent is not given"))
		}

		resizeArgs := map[string]interface{}{
			"uniq_id":        uniqIdFlag,
			"skip_fs_resize": skipFsResizeInt,
			"newsize":        configIdFlag,
		}

		// get details of existing configuration
		var cloudServerDetails apiTypes.CloudServerDetails
		if err := lwCliInst.CallLwApiInto("bleed/storm/server/details",
			map[string]interface{}{"uniq_id": uniqIdFlag},
			&cloudServerDetails); err != nil {
			lwCliInst.Die(err)
		}

		var liveResize bool
		var twoRebootResize bool
		if privateParentFlag == "" {
			// non private parent resize
			if memoryFlag != -1 || diskspaceFlag != -1 || vcpuFlag != -1 {
				lwCliInst.Die(fmt.Errorf("cannot pass --memory --diskspace or --vcpu when --private-parent is not given"))
			}

			// if already on the given config, nothing to do
			if cloudServerDetails.ConfigId == configIdFlag {
				lwCliInst.Die(fmt.Errorf("already on config_id [%d]; not initiating a resize", configIdFlag))
			}

			validateFields[configIdFlag] = "PositiveInt64"
			if err := validate.Validate(validateFields); err != nil {
				lwCliInst.Die(err)
			}

			// determine reboot expectation.
			//   resize up full: 2 reboot
			//   resize up quick (skip-fs-resize) 1 reboot
			//   resize down: 1 reboot
			var configDetails apiTypes.CloudConfigDetails
			if err := lwCliInst.CallLwApiInto("bleed/storm/config/details",
				map[string]interface{}{"id": configIdFlag}, &configDetails); err != nil {
				lwCliInst.Die(err)
			}

			if configDetails.Disk >= cloudServerDetails.DiskSpace {
				// disk space going up..
				if !skipFsResizeFlag {
					// .. and not skipping fs resize, will be 2 reboots.
					twoRebootResize = true
				}
			}
		} else {
			// private parent resize specific logic
			if memoryFlag == -1 && diskspaceFlag == -1 && vcpuFlag == -1 {
				lwCliInst.Die(fmt.Errorf(
					"resizes on private parents require at least least one of: --memory --diskspace --vcpu flags"))
			}

			privateParentUniqId, err := derivePrivateParentUniqId(privateParentFlag)
			if err != nil {
				lwCliInst.Die(err)
			}

			var (
				diskspaceChanging bool
				vcpuChanging      bool
				memoryChanging    bool
				memoryCanLive     bool
				vcpuCanLive       bool
			)
			// record what resources are changing
			if diskspaceFlag != -1 {
				if cloudServerDetails.DiskSpace != diskspaceFlag {
					diskspaceChanging = true
				}
			}
			if vcpuFlag != -1 {
				if cloudServerDetails.Vcpu != vcpuFlag {
					vcpuChanging = true
				}
			}
			if memoryFlag != -1 {
				if cloudServerDetails.Memory != memoryFlag {
					memoryChanging = true
				}
			}
			if !diskspaceChanging && !vcpuChanging && !memoryChanging {
				lwCliInst.Die(fmt.Errorf(
					"private parent resize, but passed diskspace, memory, vcpu values match existing values"))
			}

			resizeArgs["newsize"] = 0                  // 0 indicates private parent resize
			resizeArgs["parent"] = privateParentUniqId // uniq_id of the private parent
			validateFields[privateParentUniqId] = "UniqId"
			// server/resize api method always wants diskspace, vcpu, memory passed for pp resize, even if not changing
			// value. So set to current value, then override based on passed flags.
			resizeArgs["diskspace"] = cloudServerDetails.DiskSpace
			resizeArgs["memory"] = cloudServerDetails.Memory
			resizeArgs["vcpu"] = cloudServerDetails.Vcpu

			if diskspaceFlag != -1 {
				resizeArgs["diskspace"] = diskspaceFlag // desired diskspace
				validateFields[diskspaceFlag] = "PositiveInt64"
			}
			if memoryFlag != -1 {
				resizeArgs["memory"] = memoryFlag // desired memory
				validateFields[memoryFlag] = "PositiveInt64"
			}
			if vcpuFlag != -1 {
				resizeArgs["vcpu"] = vcpuFlag // desired vcpus
				validateFields[vcpuFlag] = "PositiveInt64"
			}

			// determine if this will be a live resize
			if _, exists := resizeArgs["memory"]; exists {
				if memoryFlag >= cloudServerDetails.Memory {
					// asking for more RAM
					memoryCanLive = true
				}
			}
			if _, exists := resizeArgs["vcpu"]; exists {
				if vcpuFlag >= cloudServerDetails.Vcpu {
					// asking for more vcpu
					vcpuCanLive = true
				}
			}

			if memoryFlag != -1 && vcpuFlag != -1 {
				if vcpuCanLive && memoryCanLive {
					liveResize = true
				}
			} else if memoryCanLive {
				liveResize = true
			} else if vcpuCanLive {
				liveResize = true
			}

			// if diskspace allocation changes its not currently ever done live regardless of memory, vcpu
			if diskspaceFlag != -1 {
				if resizeArgs["diskspace"] != cloudServerDetails.DiskSpace {
					liveResize = false
				}
			}
		}

		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		_, err := lwCliInst.LwApiClient.Call("bleed/server/resize", resizeArgs)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("server resized started! You can check progress with 'cloud server status --uniq_id %s'\n\n", uniqIdFlag)

		if liveResize {
			fmt.Printf("\nthis resize will be performed live without downtime.\n")
		} else {
			rebootExpectation := "one reboot"
			if twoRebootResize {
				rebootExpectation = "two reboots"
			}
			fmt.Printf(
				"\nexpect %s during this process. Your server will be online as the disk is copied to the destination.\n",
				rebootExpectation)
			if twoRebootResize {
				fmt.Printf(
					"\tTIP: Avoid the second reboot by passing --skip-fs-resize. See usage for additional details.\n")
			}
		}
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerResizeCmd)

	cloudServerResizeCmd.Flags().String("private-parent", "",
		"name or uniq_id of the private-parent. Must use when adding/removing resources to a Cloud Server on a private parent.")
	cloudServerResizeCmd.Flags().String("uniq_id", "", "uniq_id of server to resize")
	cloudServerResizeCmd.Flags().Int64("diskspace", -1, "desired diskspace (when private-parent)")
	cloudServerResizeCmd.Flags().Int64("memory", -1, "desired memory (when private-parent)")
	cloudServerResizeCmd.Flags().Bool("skip-fs-resize", false, "whether or not to skip the fs resize")
	cloudServerResizeCmd.Flags().Int64("vcpu", -1, "desired vcpu count (when private-parent)")
	cloudServerResizeCmd.Flags().Int64("config_id", -1,
		"config_id of your desired config (when !private-parent) (see 'cloud server options --configs')")

	cloudServerResizeCmd.MarkFlagRequired("uniq_id")
}
