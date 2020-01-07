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
	"os"

	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/types/api"
)

var cloudServerDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get details of a Cloud Server",
	Long: `Get details of a Cloud Server.

You can check this methods API documentation for what the returned fields mean:

https://cart.liquidweb.com/storm/api/docs/bleed/Storm/Server.html#method_details
`,
	Run: func(cmd *cobra.Command, args []string) {

		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		if uniqIdFlag == "" {
			lwCliInst.Die(fmt.Errorf("--uniq_id is a required flag"))
		}

		var details apiTypes.CloudServerDetails
		if err := lwCliInst.CallLwApiInto("bleed/storm/server/details",
			map[string]interface{}{"uniq_id": uniqIdFlag}, &details); err != nil {
			lwCliInst.Die(err)
		}

		if jsonFlag {
			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(details)
			if err != nil {
				lwCliInst.Die(err)
			}
			fmt.Printf(pretty)
			os.Exit(0)
		}

		_printCloudServerDetailsFromDetailsStruct(&details)
	},
}

func _printCloudServerDetailsFromDetailsStruct(details *apiTypes.CloudServerDetails) {
	fmt.Printf("Domain: %s UniqId: %s\n", details.Domain, details.UniqId)

	fmt.Printf("\tIp: %s\n", details.Ip)
	fmt.Printf("\tIpCount: %d\n", details.IpCount)
	fmt.Printf("\tRegion: %s (id %d) Zone: %s (id %d)\n", details.Zone.Region.Name,
		details.Zone.Region.Id, details.Zone.Name, details.Zone.Id)

	if details.PrivateParent != "" {
		fmt.Printf("\tPrivate Parent Child on Private Parent [%s]\n", details.PrivateParent)
	} else {
		fmt.Printf("\tConfigId: %d\n", details.ConfigId)
	}
	fmt.Printf("\tConfigDescription: %s\n", details.ConfigDescription)
	fmt.Printf("\tVcpus: %d\n", details.Vcpu)
	fmt.Printf("\tMemory: %d\n", details.Memory)
	fmt.Printf("\tDiskSpace: %d\n", details.DiskSpace)
	fmt.Printf("\tTemplate: %s\n", details.Template)
	fmt.Printf("\tTemplateDescription: %s\n", details.TemplateDescription)
	fmt.Printf("\tType: %s\n", details.Type)
	fmt.Printf("\tBackupEnabled: %d\n", details.BackupEnabled)
	if details.BackupEnabled == 1 {
		fmt.Printf("\tBackupPlan: %s\n", details.BackupPlan)
		fmt.Printf("\tBackupSize: %.0f\n", details.BackupSize)
		if details.BackupQuota != 0 {
			fmt.Printf("\tBackupQuota: %d\n", details.BackupQuota)
		}
	}
	fmt.Printf("\tBandwidthQuota: %s\n", details.BandwidthQuota)
	fmt.Printf("\tManageLevel: %s\n", details.ManageLevel)
	fmt.Printf("\tActive: %d\n", details.Active)
	fmt.Printf("\tAccnt: %d\n", details.Accnt)

	var attachedDetails apiTypes.CloudNetworkPrivateIsAttachedResponse
	err := lwCliInst.CallLwApiInto("bleed/network/private/isattached", map[string]interface{}{
		"uniq_id": details.UniqId}, &attachedDetails)
	if err != nil {
		lwCliInst.Die(err)
	}
	if attachedDetails.IsAttached {
		fmt.Printf("\tAttached To Private Network\n")
	}
}

func init() {
	cloudServerCmd.AddCommand(cloudServerDetailsCmd)

	cloudServerDetailsCmd.Flags().Bool("json", false, "output in json format")
	cloudServerDetailsCmd.Flags().String("uniq_id", "", "get details of this uniq_id")
}
