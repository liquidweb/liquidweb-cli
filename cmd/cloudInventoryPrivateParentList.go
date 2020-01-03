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

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/types/api"
)

var cloudInventoryPrivateParentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Private Parents on your account",
	Long:  `List Private Parents on your account`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		methodArgs := instance.AllPaginatedResultsArgs{
			Method:         "bleed/storm/private/parent/list",
			ResultsPerPage: 100,
		}
		results, err := lwCliInst.AllPaginatedResults(&methodArgs)
		if err != nil {
			lwCliInst.Die(err)
		}

		if jsonFlag {
			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(results)
			if err != nil {
				lwCliInst.Die(err)
			}
			fmt.Printf(pretty)
			os.Exit(0)
		}

		for _, item := range results.Items {
			var details apiTypes.CloudPrivateParentDetails
			if err := instance.CastFieldTypes(item, &details); err != nil {
				lwCliInst.Die(err)
			}

			_printPrivateParentDetailsFromDetailsStruct(&details)
		}
	},
}

func _printPrivateParentDetailsFromDetailsStruct(details *apiTypes.CloudPrivateParentDetails) {
	fmt.Printf("Private Parent: %s\n", details.Domain)
	fmt.Printf("\tUniqId: %s\n", details.UniqId)
	fmt.Printf("\tStatus: %s\n", details.Status)
	fmt.Printf("\tConfigId: %d\n", details.ConfigId)
	fmt.Printf("\tConfigDescription: %s\n", details.ConfigDescription)
	fmt.Printf("\tVcpus: %d\n", details.Vcpu)

	// resources
	fmt.Printf("\tResource Usage:\n")
	// diskspace
	fmt.Printf("\t\tDiskSpace:\n")
	fmt.Printf("\t\t\t%d out of %d used; free %d\n", details.Resources.DiskSpace.Used,
		details.Resources.DiskSpace.Total, details.Resources.DiskSpace.Free)
	// memory
	fmt.Printf("\t\tMemory:\n")
	fmt.Printf("\t\t\t%d out of %d used; free %d\n", details.Resources.Memory.Used,
		details.Resources.Memory.Total, details.Resources.Memory.Free)

	fmt.Printf("\tRegion %s (id %d) - %s (id %d)\n", details.Zone.Region.Name, details.Zone.Region.Id,
		details.Zone.Description, details.Zone.Id)

	fmt.Printf("\tHypervisor: %s\n", details.Zone.HvType)
	fmt.Printf("\tCreateDate: %s\n", details.CreateDate)
	fmt.Printf("\tLicenseState: %s\n", details.LicenseState)
}

func init() {
	cloudInventoryPrivateParentCmd.AddCommand(cloudInventoryPrivateParentListCmd)

	cloudInventoryPrivateParentListCmd.Flags().Bool("json", false, "output in json format")
}
