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
	"strings"

	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/types/api"
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

If this is a resize of a private parent child instance, pass --private-parent
with a value of the name of the private parent. When passing --private-parent,
the following flags are required:

  --diskspace
  --memory
  --vcpu

If you resize a private parent child instance, and only up the memory or vcpu,
this will be applied live without downtime to your Cloud Server.
`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		diskspaceFlag, _ := cmd.Flags().GetInt64("diskspace")
		configIdFlag, _ := cmd.Flags().GetInt64("config_id")
		memoryFlag, _ := cmd.Flags().GetInt64("memory")
		parentFlag, _ := cmd.Flags().GetString("parent")
		skipFsResizeFlag, _ := cmd.Flags().GetBool("skip-fs-resize")
		vcpuFlag, _ := cmd.Flags().GetInt64("vcpu")
		privateParentFlag, _ := cmd.Flags().GetString("private-parent")

		// convert bool to int for api
		skipFsResizeInt := 0
		if skipFsResizeFlag {
			skipFsResizeInt = 1
		}

		if uniqIdFlag == "" {
			lwCliInst.Die(fmt.Errorf("flag --uniq_id is required"))
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
		if err := lwCliInst.CallLwApiInto("bleed/storm/server/details", map[string]interface{}{"uniq_id": uniqIdFlag}, &cloudServerDetails); err != nil {
			lwCliInst.Die(err)
		}

		var liveResize bool
		var oneRebootResize bool
		if privateParentFlag == "" {
			// non private parent resize
			if memoryFlag != -1 || diskspaceFlag != -1 || vcpuFlag != -1 {
				lwCliInst.Die(fmt.Errorf("cannot pass --memory --diskspace or --vcpu when --private-parent is not given"))
			}

			// determine reboot expectation.
		} else {
			// private parent resize specific logic
			if memoryFlag == -1 && diskspaceFlag == -1 && vcpuFlag == -1 {
				lwCliInst.Die(fmt.Errorf("resizes on private parents require at least least one of: --memory --diskspace --vcpu flags"))
			}

			var privateParentUniqId string
			var privateParentDetails apiTypes.CloudPrivateParentDetails
			var privateParentDetailsErr error

			// if privateParentFlag looks like a uniq_id, try it as a uniq_id first.
			if len(privateParentFlag) == 6 && strings.ToUpper(privateParentFlag) == privateParentFlag {
				if err := lwCliInst.CallLwApiInto("bleed/storm/private/parent/details", map[string]interface{}{"uniq_id": privateParentFlag}, &privateParentDetails); err == nil {
					privateParentUniqId = privateParentFlag
				} else {
					privateParentDetailsErr = errors.New(fmt.Sprintf("failed fetching parent details treating given --private-parent arg as a uniq_id [%s]: %s", privateParentFlag, err))
				}
			}

			// if we havent found the pp details yet, try assuming privateParentFlag is the name of the pp
			if privateParentUniqId == "" {
				methodArgs := instance.AllPaginatedResultsArgs{
					Method:         "bleed/storm/private/parent/list",
					ResultsPerPage: 100,
				}
				results, err := lwCliInst.AllPaginatedResults(&methodArgs)
				if err != nil {
					lwCliInst.Die(err)
				}

				for _, item := range results.Items {
					var privateParentDetails apiTypes.CloudPrivateParentDetails
					if err := instance.CastFieldTypes(item, &privateParentDetails); err != nil {
						lwCliInst.Die(err)
					}

					if privateParentDetails.Domain == privateParentFlag {
						// found it get details
						err := lwCliInst.CallLwApiInto("bleed/storm/private/parent/details",
							map[string]interface{}{
								"uniq_id": privateParentDetails.UniqId,
							},
							&privateParentDetails)
						if err != nil {
							privateParentDetailsErr = fmt.Errorf("failed fetching private parent details for discovered uniq_id [%s] error: %s %w",
								privateParentDetails.UniqId, err, privateParentDetailsErr)
							lwCliInst.Die(privateParentDetailsErr)
						}
						privateParentUniqId = privateParentDetails.UniqId
						break // found the uniq_id so break
					}
				}
			}

			if privateParentUniqId == "" {
				lwCliInst.Die(fmt.Errorf("failed deriving uniq_id from --private-parent [%s]: %s", privateParentFlag, privateParentDetailsErr))
			}

			var (
				diskspaceChanging bool
				vcpuChanging      bool
				memoryChanging    bool
			)
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
					"private parent resize, but your passed diskspace, memory, vcpu values match what the Cloud Server already has.. no need to resize"))
			}

			resizeArgs["newsize"] = 0                  // 0 indicates private parent resize
			resizeArgs["parent"] = privateParentUniqId // uniq_id of the private parent
			// server/resize api method always wants diskspace, vcpu, memory passed for pp resize, even if not changing
			// value. So set to current value, then override based on passed flags.
			resizeArgs["diskspace"] = cloudServerDetails.DiskSpace
			resizeArgs["memory"] = cloudServerDetails.Memory
			resizeArgs["vcpu"] = cloudServerDetails.Vcpu

			if diskspaceFlag != -1 {
				resizeArgs["diskspace"] = diskspaceFlag // desired diskspace
			}
			if memoryFlag != -1 {
				resizeArgs["memory"] = memoryFlag // desired memory
			}
			if vcpuFlag != -1 {
				resizeArgs["vcpu"] = vcpuFlag // desired vcpus
			}
			if parentFlag != "" {
				resizeArgs["parent"] = privateParentFlag // name of the private parent
			}

			// determine if resize can be performed live.
			var (
				memoryCanLive bool
				vcpuCanLive   bool
			)

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

		_, err := lwCliInst.LwApiClient.Call("bleed/server/resize", resizeArgs)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("server resized started! You can check progress with 'cloud server status --uniq_id %s'\n\n", uniqIdFlag)

		if liveResize {
			fmt.Printf("\nthis resize will be performed live without downtime.\n")
		} else {
			rebootExpectation := "two"
			if oneRebootResize {
				rebootExpectation = "one"
			}
			fmt.Printf("\nexpect %s reboot during this process. Your server will be online as the disk is copied to the destination.\n", rebootExpectation)
		}
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerResizeCmd)

	cloudServerResizeCmd.Flags().String("private-parent", "", "name or uniq_id of the private-parent. Must use when adding/removing resources to a Cloud Server on a private parent.")
	cloudServerResizeCmd.Flags().String("uniq_id", "", "uniq_id of server to resize")
	cloudServerResizeCmd.Flags().Int64("diskspace", -1, "desired diskspace (required when private-parent)")
	cloudServerResizeCmd.Flags().Int64("memory", -1, "desired memory (required when private-parent)")
	cloudServerResizeCmd.Flags().String("parent", "", "name of private parent (required when private-parent)")
	cloudServerResizeCmd.Flags().Bool("skip-fs-resize", false, "whether or not to skip the fs resize")
	cloudServerResizeCmd.Flags().Int64("vcpu", -1, "desired vcpu count (required when private-parent)")
	cloudServerResizeCmd.Flags().Int64("config_id", -1, "config_id of your desired config (dont use with private-parent) (see 'cloud server options --configs')")
}
