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
)

var cloudServerResizeCmd = &cobra.Command{
	Use:   "resize",
	Short: "Resize a Cloud Server",
	Long: `Resize a Cloud Server.

Resize a Cloud Server to a new config. Available config-id's can be found in
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
uniq-id. When passing --private-parent, at least one of the following flags
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
		params := &instance.CloudServerResizeParams{}

		params.UniqId, _ = cmd.Flags().GetString("uniq-id")
		params.DiskSpace, _ = cmd.Flags().GetInt64("diskspace")
		params.ConfigId, _ = cmd.Flags().GetInt64("config-id")
		params.Memory, _ = cmd.Flags().GetInt64("memory")
		params.SkipFsResize, _ = cmd.Flags().GetBool("skip-fs-resize")
		params.Vcpu, _ = cmd.Flags().GetInt64("vcpu")
		params.PrivateParent, _ = cmd.Flags().GetString("private-parent")

		status, err := lwCliInst.CloudServerResize(params)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(status)
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerResizeCmd)

	cloudServerResizeCmd.Flags().String("private-parent", "",
		"name or uniq-id of the private-parent. Must use when adding/removing resources to a Cloud Server on a private parent.")
	cloudServerResizeCmd.Flags().String("uniq-id", "", "uniq-id of server to resize")
	cloudServerResizeCmd.Flags().Int64("diskspace", -1, "desired diskspace (when private-parent)")
	cloudServerResizeCmd.Flags().Int64("memory", -1, "desired memory (when private-parent)")
	cloudServerResizeCmd.Flags().Bool("skip-fs-resize", false, "whether or not to skip the fs resize")
	cloudServerResizeCmd.Flags().Int64("vcpu", -1, "desired vcpu count (when private-parent)")
	cloudServerResizeCmd.Flags().Int64("config-id", -1,
		"config-id of your desired config (when !private-parent) (see 'cloud server options --configs')")

	if err := cloudServerResizeCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
