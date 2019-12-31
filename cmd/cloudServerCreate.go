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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/instance"
)

var cloudServerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Cloud Server",
	Long: `Create a Cloud Server.

Requires various flags. Please see the flag section of help.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bleed/server/create")
		/*
			{
				"params": {
					"domain": "test.create-this-is-junk.io",
					"type": "SS.VPS",
					"features": {
						"ConfigId": "41",
						"Template": "CENTOS_8_UNMANAGED",
						"Bandwidth": "SS.5000",
						"ExtraIp": {
							"value": 1,
							"count": 0
						},
						"LiquidWebBackupPlan": "Quota",
						"BackupQuota": "50"
					},
					"pool_ips": [],
					"new_ips": 1,
					"public_ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDJ4ROyFN7EObvHSalBZbNiRU/o9jR7caBhs2Vxwpaio1qul9Z8YGM8O4DgPOyc+6VTORsYkth3ylRMYvtPrJyE8Rd3NlTqYuAoxKR14S2dNGcJxOmLIZ12kTOW8dqlczgmRQNEH7uDTV22rR/O9p4SdYOgMZndyrL8R0w1Ci96l8LWOK7NYcjlEBtH079f96sKar41lLY7CDf2jq6cxuCpONt+tAtOIpClzeF7b30HpJszlwkypgmUWwYZUPIy/sEY3uXmRtRXMdzDhQKHmCXmtAclvbCGMH/r1ilH2A5lz4nw3YFctxYP1sWQ6tvlwiqxEMIyuYGW8He5cj7zPUhP4bSjf1qC+67o18PzX/YzOWxKPNicde/B9xdYGCxX+SGWRIfdn4VVEbe358xTTJuJZhNkNsClLa7tgqqQS19Ai8Cw6PrSn6ow52PSSXz0mwwGxObYBe29QwziPwazv1wLBpVLAWRoQbnNWvvfIF5rTZxctEROONghofdyj5c3BgJyjVXMjYbVHQqhMsu4SJ6DP+M6PZwXJgyKHDdGSyH+i2Mqc2NWXYVJXgafm2q9vWphoXxW+XoMFjylc3DXuCr4L/qvd9yXvp0GRJVc2Fyd7bSO9ND7m6SYK5Dew8uLF9kBdXiNOrD3xeNgnMyHqLp0eEjEH2bb7+kUfuZKI/okZw== ssullivan@data.wks.liquidweb.com",
					"zone": "86172"
				}
			}
		*/

	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerCreateCmd)

	var sshPubKeyFile string
	home, err := homedir.Dir()
	if err == nil {
		sshPubKeyFile = fmt.Sprintf("%s/.ssh/id_rsa.pub", home)
	}

	defaultHostname := fmt.Sprintf("%s.%s.io", instance.RandomString(4), instance.RandomString(10))

	cloudServerCreateCmd.Flags().String("template", "", "name of the template to use")
	cloudServerCreateCmd.Flags().String("type", "SS.VPS", "some examples of types; SS.VPS, SS.VPS.WIN, SS.VM, SS.VM.WIN")
	cloudServerCreateCmd.Flags().String("domain", defaultHostname, "hostname to set")
	cloudServerCreateCmd.Flags().Int("ips", 1, "amount of IP addresses")
	// TODO pool_ips
	cloudServerCreateCmd.Flags().String("public-ssh-key", sshPubKeyFile, "path to file containing the public ssh key you wish to be on the new cloud server")
	cloudServerCreateCmd.Flags().Int("config-id", 0, "config_id to use")
	cloudServerCreateCmd.Flags().String("backup-plan", "None", "LiquidWeb cloud server backup plan to use")
	cloudServerCreateCmd.Flags().Int("backup-plan-quota", 300, "Quota amount. Should only be used with '--backup-plan Quota'")
	cloudServerCreateCmd.Flags().String("bandwidth", "SS.1000", "bandwidth package to use")
	cloudServerCreateCmd.Flags().Int("zone", 0, "Cloud server zone to create in")
}
