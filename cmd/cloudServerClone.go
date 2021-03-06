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

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/utils"
	"github.com/liquidweb/liquidweb-cli/validate"
)

var cloudServerCloneCmdPoolIpsFlag []string
var cloudServerCloneCmdPool6IpsFlag []string

var cloudServerCloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a Cloud Server",
	Long: `Clone a Cloud Server.

Clones an existing Cloud Server to a new Cloud Server.

During the clone operation there will be no downtime for the source Cloud Server.

Any optional omitted flag will be defaulted to the value from the source Cloud
Server where possible.

** Cloning to a Private Parent: **

--private-parent must be passed containing either the uniq-id of the Private Parent,
or its name.

The flags --diskspace --vcpu --memory must all be passed if the source Cloud
Server is not on a Private Parent.`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")
		passwordFlag, _ := cmd.Flags().GetString("password")
		zoneFlag, _ := cmd.Flags().GetInt64("zone")
		newIpsFlag, _ := cmd.Flags().GetInt64("new-ips")
		newIp6sFlag, _ := cmd.Flags().GetInt64("new-ip6s")
		hostnameFlag, _ := cmd.Flags().GetString("hostname")
		privateParentFlag, _ := cmd.Flags().GetString("private-parent")
		diskspaceFlag, _ := cmd.Flags().GetInt64("diskspace")
		memoryFlag, _ := cmd.Flags().GetInt64("memory")
		vcpuFlag, _ := cmd.Flags().GetInt64("vcpu")
		configIdFlag, _ := cmd.Flags().GetInt64("config-id")

		validateFields := map[interface{}]interface{}{
			uniqIdFlag: "UniqId",
			// expanded out struct to show ability.. its treated as required like above
			hostnameFlag: map[string]string{"type": "NonEmptyString", "optional": "false"},
		}

		if privateParentFlag == "" && configIdFlag == -1 {
			lwCliInst.Die(fmt.Errorf("must pass --config-id or --private-parent"))
		}

		var privateParentUniqId string
		if privateParentFlag != "" {
			var (
				err  error
				zone int64
			)
			privateParentUniqId, zone, err = lwCliInst.DerivePrivateParentUniqId(privateParentFlag)
			if err != nil {
				lwCliInst.Die(err)
			}
			if zoneFlag == -1 {
				zoneFlag = zone
			}
		}

		// buildout api bleed/server/clone parameters
		cloneArgs := map[string]interface{}{
			"uniq_id":  uniqIdFlag,
			"domain":   hostnameFlag,
			"new_ips":  newIpsFlag,
			"new_ip6s": newIp6sFlag,
		}
		if passwordFlag != "" {
			cloneArgs["password"] = passwordFlag
		}
		if zoneFlag != -1 {
			cloneArgs["zone"] = zoneFlag
			validateFields[zoneFlag] = "PositiveInt64"
		}
		if privateParentUniqId != "" {
			cloneArgs["parent"] = privateParentUniqId
		}
		if diskspaceFlag != -1 {
			cloneArgs["diskspace"] = diskspaceFlag
			validateFields[diskspaceFlag] = "PositiveInt64"
		}
		if memoryFlag != -1 {
			cloneArgs["memory"] = memoryFlag
			validateFields[memoryFlag] = "PositiveInt64"
		}
		if vcpuFlag != -1 {
			cloneArgs["vcpu"] = vcpuFlag
			validateFields[vcpuFlag] = "PositiveInt64"
		}
		if configIdFlag != -1 && privateParentFlag == "" {
			cloneArgs["config_id"] = configIdFlag
			validateFields[configIdFlag] = "PositiveInt64"
		}
		if len(cloudServerCloneCmdPoolIpsFlag) > 0 {
			cloneArgs["pool_ips"] = cloudServerCloneCmdPoolIpsFlag
			for _, ip := range cloudServerCloneCmdPoolIpsFlag {
				validateFields[ip] = "IP"
			}
		}
		if len(cloudServerCloneCmdPool6IpsFlag) > 0 {
			cloneArgs["pool6_ips"] = cloudServerCloneCmdPool6IpsFlag
			for _, ip := range cloudServerCloneCmdPool6IpsFlag {
				validateFields[ip] = "CIDR"
			}
		}

		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		var details apiTypes.CloudServerCloneResponse
		if err := lwCliInst.CallLwApiInto("bleed/server/clone", cloneArgs, &details); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf(
			"Success! Cloning existing Cloud Server [%s] to new Cloud Server [%s]. Check status with 'cloud server status --uniq-id %s'\n",
			uniqIdFlag, details.UniqId, uniqIdFlag)

	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerCloneCmd)

	// General
	cloudServerCloneCmd.Flags().String("uniq-id", "", "uniq-id of Cloud Server to clone")
	cloudServerCloneCmd.Flags().String("password", "", "root or administrator password for new Cloud Server")
	cloudServerCloneCmd.Flags().Int64("zone", -1, "zone for new Cloud Server")
	cloudServerCloneCmd.Flags().Int64("new-ips", 1, "amount of IPv4 addresses for the new Cloud Server")
	cloudServerCloneCmd.Flags().Int64("new-ip6s", 0, "amount of new IPv6 /64's for the new Cloud Server")
	cloudServerCloneCmd.Flags().String("hostname", fmt.Sprintf("%s.%s.io", utils.RandomString(4),
		utils.RandomString(10)), "hostname for new Cloud Server")
	cloudServerCloneCmd.Flags().StringSliceVar(&cloudServerCloneCmdPoolIpsFlag, "pool-ips", []string{},
		"IPv4 ips from your IP Pool separated by ',' to assign to the new Cloud Server")
	cloudServerCloneCmd.Flags().StringSliceVar(&cloudServerCloneCmdPool6IpsFlag, "pool6-ips", []string{},
		"IPv6 ips from your IP Pool separated by ',' to assign to the new Cloud Server")

	// Private Parent
	cloudServerCloneCmd.Flags().String("private-parent", "",
		"name or uniq-id of the Private Parent to place new Cloud Server on (see: 'cloud private-parent list')")
	cloudServerCloneCmd.Flags().Int64("diskspace", -1, "diskspace for new Cloud Server (when private-parent)")
	cloudServerCloneCmd.Flags().Int64("memory", -1, "memory for new Cloud Server (when private-parent)")
	cloudServerCloneCmd.Flags().Int64("vcpu", -1, "amount of vcpus for new Cloud Server (when private-parent)")

	// Non Private Parent
	cloudServerCloneCmd.Flags().Int64("config-id", cast.ToInt64(defaultFlag("cloud_server_clone_config-id", -1)),
		"config-id for new Cloud Server (when !private-parent) (see: 'cloud server options --configs')")

	if err := cloudServerCloneCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
