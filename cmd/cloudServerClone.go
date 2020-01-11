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
	"github.com/liquidweb/liquidweb-cli/utils"
)

var cloudServerCloneCmdPoolIpsFlag []string

var cloudServerCloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a Cloud Server",
	Long: `Clone a Cloud Server.

Clones an existing Cloud Server to a new Cloud Server.

During the clone operation there will be no downtime for the source Cloud Server.

Any optional omitted flag will be defaulted to the value from the source Cloud
Server where possible.

** Cloning to a Private Parent: **

--private-parent must be passed containg either the uniq_id of the Private Parent,
or its name.

The flags --diskspace --vcpu --memory must all be passed if the source Cloud
Server is not on a Private Parent.`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		passwordFlag, _ := cmd.Flags().GetString("password")
		zoneFlag, _ := cmd.Flags().GetInt64("zone")
		newIpsFlag, _ := cmd.Flags().GetInt64("new_ips")
		hostnameFlag, _ := cmd.Flags().GetString("hostname")
		privateParentFlag, _ := cmd.Flags().GetString("private-parent")
		diskspaceFlag, _ := cmd.Flags().GetInt64("diskspace")
		memoryFlag, _ := cmd.Flags().GetInt64("memory")
		vcpuFlag, _ := cmd.Flags().GetInt64("vcpu")
		configIdFlag, _ := cmd.Flags().GetInt64("config_id")

		// flag check
		if uniqIdFlag == "" {
			lwCliInst.Die(fmt.Errorf("Flag --uniq_id is required"))
		}

		if privateParentFlag != "" && configIdFlag != -1 {
			lwCliInst.Die(fmt.Errorf("cant pass both --config_id and --private-parent flags"))
		}

		var privateParentUniqId string
		if privateParentFlag != "" {
			var err error
			privateParentUniqId, err = derivePrivateParentUniqId(privateParentFlag)
			if err != nil {
				lwCliInst.Die(err)
			}
		}

		// buildout api bleed/server/clone parameters
		cloneArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
			"domain":  hostnameFlag,
			"new_ips": newIpsFlag,
		}
		if passwordFlag != "" {
			cloneArgs["password"] = passwordFlag
		}
		if zoneFlag != -1 {
			cloneArgs["zone"] = zoneFlag
		}
		if privateParentUniqId != "" {
			cloneArgs["parent"] = privateParentUniqId
		}
		if diskspaceFlag != -1 {
			cloneArgs["diskspace"] = diskspaceFlag
		}
		if memoryFlag != -1 {
			cloneArgs["memory"] = memoryFlag
		}
		if vcpuFlag != -1 {
			cloneArgs["vcpu"] = vcpuFlag
		}
		if configIdFlag != -1 {
			cloneArgs["config_id"] = configIdFlag
		}
		if len(cloudServerCloneCmdPoolIpsFlag) > 0 {
			cloneArgs["pool_ips"] = cloudServerCloneCmdPoolIpsFlag
		}

		var details apiTypes.CloudServerCloneResponse
		if err := lwCliInst.CallLwApiInto("bleed/server/clone", cloneArgs, &details); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf(
			"Success! Cloning existing Cloud Server [%s] to new Cloud Server [%s]. Check status with 'cloud server status --uniq_id %s'\n",
			uniqIdFlag, details.UniqId, uniqIdFlag)

	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerCloneCmd)

	// General
	cloudServerCloneCmd.Flags().String("uniq_id", "", "uniq_id of Cloud Server to clone")
	cloudServerCloneCmd.Flags().String("password", "", "root or administrator password for new Cloud Server")
	cloudServerCloneCmd.Flags().Int64("zone", -1, "zone for new Cloud Server")
	cloudServerCloneCmd.Flags().Int64("new_ips", 1, "amount of IP addresses for new Cloud Server")
	cloudServerCloneCmd.Flags().String("hostname", fmt.Sprintf("%s.%s.io", utils.RandomString(4),
		utils.RandomString(10)), "hostname for new Cloud Server")
	cloudServerCloneCmd.Flags().StringSliceVar(&cloudServerCloneCmdPoolIpsFlag, "pool-ips", []string{},
		"ips from your IP Pool separated by ',' to assign to the new Cloud Server")

	// Private Parent
	cloudServerCloneCmd.Flags().String("private-parent", "",
		"name or uniq_id of the Private Parent to place new Cloud Server on (see: 'cloud inventory private-parent list')")
	cloudServerCloneCmd.Flags().Int64("diskspace", -1, "diskspace for new Cloud Server (when private-parent)")
	cloudServerCloneCmd.Flags().Int64("memory", -1, "memory for new Cloud Server (when private-parent)")
	cloudServerCloneCmd.Flags().Int64("vcpu", -1, "amount of vcpus for new Cloud Server (when private-parent)")

	// Non Private Parent
	cloudServerCloneCmd.Flags().Int64("config_id", -1,
		"config_id for new Cloud Server (when !private-parent) (see: 'cloud server options --configs')")
}
