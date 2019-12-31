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

If this is a resize of a private parent child instance, pass --newsize with
a 0 value. When newsize is 0, the following flags are required:

  --diskspace
  --memory
  --parent
  --vcpu

If you resize a private parent child instance, and only up the memory or vcpu,
this will be applied live without downtime to your Cloud Server.
`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		jsonFlag, _ := cmd.Flags().GetBool("json")
		verboseFlag, _ := cmd.Flags().GetBool("verbose")
		diskspaceFlag, _ := cmd.Flags().GetInt("diskspace")
		newsizeFlag, _ := cmd.Flags().GetInt("newsize")
		memoryFlag, _ := cmd.Flags().GetInt("memory")
		parentFlag, _ := cmd.Flags().GetString("parent")
		skipFsResizeFlag, _ := cmd.Flags().GetBool("skip-fs-resize")
		vcpuFlag, _ := cmd.Flags().GetInt("vcpu")

		// convert bool to int for api
		skipFsResizeInt := 0
		if skipFsResizeFlag {
			skipFsResizeInt = 1
		}

		if uniqIdFlag == "" {
			lwCliInst.Die(fmt.Errorf("flag --uniq_id is required"))
		}
		if newsizeFlag == -1 {
			lwCliInst.Die(fmt.Errorf("flag --newsize is required"))
		}

		resizeArgs := map[string]interface{}{
			"uniq_id":        uniqIdFlag,
			"skip_fs_resize": skipFsResizeInt,
			"newsize":        newsizeFlag,
		}

		if newsizeFlag == 0 {
			// private parent resize, dtrt with private parent args
			if memoryFlag == -1 && diskspaceFlag == -1 && vcpuFlag == -1 || parentFlag == "" {
				lwCliInst.Die(fmt.Errorf("when --newsize is 0, the parent flag is required. Also at least one of memory, diskspace, vcpu flags must be passed."))
			}
			resizeArgs["diskspace"] = diskspaceFlag
			resizeArgs["memory"] = memoryFlag
			resizeArgs["vcpu"] = vcpuFlag
			resizeArgs["parent"] = parentFlag
		}

		if verboseFlag {
			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(resizeArgs)
			if err == nil {
				fmt.Printf("sending resize args:\n")
				fmt.Printf(pretty)
			}
		}

		result, err := lwCliInst.LwApiClient.Call("bleed/server/resize", resizeArgs)
		if err != nil {
			lwCliInst.Die(err)
		}

		if jsonFlag {
			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(result)
			if err != nil {
				lwCliInst.Die(err)
			}
			fmt.Printf(pretty)
		} else {
			fmt.Println("server resized started!\n\n")
			// TODO update downtime expectation text once logic is added to detect if PP child will have ram/cpu hotplug performed
			//fmt.Printf("Expect one reboot during this process. Your server will be online as the disk is copied to the destination.\n")
		}
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerResizeCmd)

	cloudServerResizeCmd.Flags().Bool("verbose", false, "be verbose")
	cloudServerResizeCmd.Flags().Bool("json", false, "output in json format")
	cloudServerResizeCmd.Flags().String("uniq_id", "", "uniq_id of server to resize")
	cloudServerResizeCmd.Flags().Int("diskspace", -1, "desired diskspace (when --newsize is 0)")
	cloudServerResizeCmd.Flags().Int("memory", -1, "desired memory (when --newsize is 0)")
	cloudServerResizeCmd.Flags().String("parent", "", "name of private parent (when --newsize is 0)")
	cloudServerResizeCmd.Flags().Bool("skip-fs-resize", false, "whether or not to skip the fs resize")
	cloudServerResizeCmd.Flags().Int("vcpu", -1, "desired vcpu count (when --newsize is 0)")
	cloudServerResizeCmd.Flags().Int("newsize", -1, "should be a config_id of a config or 0 for a private parent resize")
}
