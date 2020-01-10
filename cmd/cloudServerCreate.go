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
	"io/ioutil"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/types/api"
)

var cloudServerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Cloud Server",
	Long: `Create a Cloud Server.

Requires various flags. Please see the flag section of help.

Examples:

# Create a Cloud Server on a Private Parent named "private"
'cloud server create --private-parent private --memory 1024 --diskspace 40 --vcpu 2 --zone 40460 --template DEBIAN_10_UNMANAGED'

# Create a Cloud Server on config_id 1
'cloud server create --config_id 1 --template DEBIAN_10_UNMANAGED --zone 40460'

# Create a Cloud Server from image id 111
'cloud server create --image-id 111 --zone 40460 --config_id 1'

# Create a Cloud Server from backup id 111
'cloud server create --backup-id 111 --zone 40460 --config_id 1'

These examples use default values for various flags, such as password, type, ssh-key, hostname, etc.

For a list of Templates, Configs, and Region/Zones, see 'cloud server options --configs --templates --zones'
For a list of images, see 'cloud inventory images list'
For a list of backups, see 'cloud inventory backups list'
`,
	Run: func(cmd *cobra.Command, args []string) {
		templateFlag, _ := cmd.Flags().GetString("template")
		typeFlag, _ := cmd.Flags().GetString("type")
		hostnameFlag, _ := cmd.Flags().GetString("hostname")
		ipsFlag, _ := cmd.Flags().GetInt("ips")
		pubSshKeyFlag, _ := cmd.Flags().GetString("public-ssh-key")
		configIdFlag, _ := cmd.Flags().GetInt("config_id")
		backupPlanFlag, _ := cmd.Flags().GetString("backup-plan")
		backupPlanQuotaFlag, _ := cmd.Flags().GetInt("backup-plan-quota")
		bandwidthFlag, _ := cmd.Flags().GetString("bandwidth")
		zoneFlag, _ := cmd.Flags().GetInt("zone")
		winavFlag, _ := cmd.Flags().GetString("winav")
		msSqlFlag, _ := cmd.Flags().GetString("ms_sql")
		privateParentFlag, _ := cmd.Flags().GetString("private-parent")
		passwordFlag, _ := cmd.Flags().GetString("password")
		memoryFlag, _ := cmd.Flags().GetInt("memory")
		diskspaceFlag, _ := cmd.Flags().GetInt("diskspace")
		vcpuFlag, _ := cmd.Flags().GetInt("vcpu")
		backupIdFlag, _ := cmd.Flags().GetInt("backup-id")
		imageIdFlag, _ := cmd.Flags().GetInt("image-id")

		// sanity check flags
		if zoneFlag == 0 {
			lwCliInst.Die(fmt.Errorf("--zone is a required flag"))
		}
		if configIdFlag == 0 && privateParentFlag == "" {
			lwCliInst.Die(fmt.Errorf("--config_id is a required flag without --private-parent"))
		}
		if templateFlag == "" && backupIdFlag == -1 && imageIdFlag == -1 {
			lwCliInst.Die(fmt.Errorf("at least one of the following flags must be set --template --image-id --backup-id"))
		}

		var publicSshKeyContents string
		sshPkeyContents, err := ioutil.ReadFile(pubSshKeyFlag)
		if err == nil {
			publicSshKeyContents = cast.ToString(sshPkeyContents)
		}

		// if passed a private-parent flag, derive its uniq_id
		var privateParentUniqId string
		if privateParentFlag != "" {
			privateParentUniqId, err = derivePrivateParentUniqId(privateParentFlag)
			if err != nil {
				lwCliInst.Die(err)
			}
		}

		// buildout args for bleed/server/create
		createArgs := map[string]interface{}{
			"domain":   hostnameFlag,
			"type":     typeFlag,
			"pool_ips": []string{}, // TODO
			"new_ips":  ipsFlag,
			"zone":     zoneFlag,
			"password": passwordFlag,
			"features": map[string]interface{}{
				"Bandwidth": bandwidthFlag,
				"ConfigId":  configIdFlag,
				"ExtraIp": map[string]interface{}{
					"value": ipsFlag,
					"count": 0,
				},
				"LiquidWebBackupPlan": backupPlanFlag,
			},
		}

		var isWindows bool
		if templateFlag != "" {
			createArgs["features"].(map[string]interface{})["Template"] = templateFlag
			if strings.Contains(strings.ToUpper(templateFlag), "WINDOWS") {
				isWindows = true
			}
		}
		if backupIdFlag != -1 {
			// check backup and see if its windows
			apiArgs := map[string]interface{}{"id": backupIdFlag}
			var details apiTypes.CloudBackupDetails
			err := lwCliInst.CallLwApiInto("bleed/storm/backup/details", apiArgs, &details)
			if err != nil {
				lwCliInst.Die(err)
			}
			if strings.Contains(strings.ToUpper(details.Template), "WINDOWS") {
				isWindows = true
			}
			createArgs["backup_id"] = backupIdFlag
		}
		if imageIdFlag != -1 {
			// check image and see if its windows
			apiArgs := map[string]interface{}{"id": imageIdFlag}
			var details apiTypes.CloudImageDetails
			err := lwCliInst.CallLwApiInto("bleed/storm/image/details", apiArgs, &details)
			if err != nil {
				lwCliInst.Die(err)
			}
			if strings.Contains(strings.ToUpper(details.Template), "WINDOWS") {
				isWindows = true
			}
			createArgs["image_id"] = imageIdFlag
		}

		// windows servers need special arguments
		if isWindows {
			if winavFlag == "" {
				winavFlag = "None"
			}
			createArgs["features"].(map[string]interface{})["WinAV"] = winavFlag
			createArgs["features"].(map[string]interface{})["WindowsLicense"] = "Windows"
			if typeFlag == "SS.VPS" {
				createArgs["type"] = "SS.VPS.WIN"
			}
			if msSqlFlag == "" {
				msSqlFlag = "None"
			}
			var coreCnt int
			if vcpuFlag == -1 {
				// standard config_id create, fetch configs core count and use it
				var details apiTypes.CloudConfigDetails
				if err := lwCliInst.CallLwApiInto("bleed/storm/config/details",
					map[string]interface{}{"id": configIdFlag}, &details); err != nil {
					lwCliInst.Die(err)
				}
				coreCnt = cast.ToInt(details.Vcpu)
			} else {
				// private parent, use vcpu flag
				coreCnt = vcpuFlag
			}
			createArgs["features"].(map[string]interface{})["MsSQL"] = map[string]interface{}{
				"value": msSqlFlag,
				"count": coreCnt,
			}
		}

		if privateParentUniqId != "" {
			// create on a private parent. diskspace, memory, vcpu are now required.
			if memoryFlag == -1 || diskspaceFlag == -1 || vcpuFlag == -1 {
				lwCliInst.Die(fmt.Errorf("flags --diskspace --memory --vcpu are required when --private-parent is passed"))
			}

			createArgs["parent"] = privateParentUniqId
			createArgs["vcpu"] = vcpuFlag
			createArgs["diskspace"] = diskspaceFlag
			createArgs["memory"] = memoryFlag
		}

		if backupPlanFlag == "Quota" {
			createArgs["features"].(map[string]interface{})["BackupQuota"] = backupPlanQuotaFlag
		}

		if publicSshKeyContents != "" {
			createArgs["public_ssh_key"] = publicSshKeyContents
		}

		result, err := lwCliInst.LwApiClient.Call("bleed/server/create", createArgs)
		if err != nil {
			lwCliInst.Die(err)
		}

		resultUniqId := result.(map[string]interface{})["uniq_id"]
		fmt.Printf(
			"Cloud server with uniq_id [%s] creating. Check status with 'cloud server status --uniq_id %s'\n",
			resultUniqId, resultUniqId)
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerCreateCmd)

	var sshPubKeyFile string
	home, err := homedir.Dir()
	if err == nil {
		sshPubKeyFile = fmt.Sprintf("%s/.ssh/id_rsa.pub", home)
	}

	randomHostname := fmt.Sprintf("%s.%s.io", instance.RandomString(4), instance.RandomString(10))
	randomPassword := instance.RandomString(25)

	cloudServerCreateCmd.Flags().String("template", "", "name of the template to use")
	cloudServerCreateCmd.Flags().String("type", "SS.VPS", "some examples of types; SS.VPS, SS.VPS.WIN, SS.VM, SS.VM.WIN")
	cloudServerCreateCmd.Flags().String("hostname", randomHostname, "hostname to set")
	cloudServerCreateCmd.Flags().Int("ips", 1, "amount of IP addresses")
	cloudServerCreateCmd.Flags().String("public-ssh-key", sshPubKeyFile,
		"path to file containing the public ssh key you wish to be on the new Cloud Server")
	cloudServerCreateCmd.Flags().Int("config_id", 0, "config_id to use")
	cloudServerCreateCmd.Flags().String("backup-plan", "None", "Cloud Server backup plan to use")
	cloudServerCreateCmd.Flags().Int("backup-plan-quota", 300, "Quota amount. Should only be used with '--backup-plan Quota'")
	cloudServerCreateCmd.Flags().String("bandwidth", "SS.10000", "bandwidth package to use")
	cloudServerCreateCmd.Flags().Int("zone", 0, "zone to create in")
	cloudServerCreateCmd.Flags().String("password", randomPassword, "root or administrator password to set")

	cloudServerCreateCmd.Flags().Int("backup-id", -1, "id of backup to create from (see 'cloud inventory backup list')")
	cloudServerCreateCmd.Flags().Int("image-id", -1, "id of image to create from (see 'cloud inventory image list')")

	// private parent specific
	cloudServerCreateCmd.Flags().String("private-parent", "",
		"name or uniq_id of the private-parent. Must use when creating a Cloud Server on a private parent.")
	cloudServerCreateCmd.Flags().Int("memory", -1, "memory (ram) value use with --private-parent")
	cloudServerCreateCmd.Flags().Int("diskspace", -1, "diskspace value use with --private-parent")
	cloudServerCreateCmd.Flags().Int("vcpu", -1, "vcpu value use with --private-parent")

	// windows specific
	cloudServerCreateCmd.Flags().String("winav", "", "Use only with Windows Servers. Typically (None or NOD32) for value when set")
	cloudServerCreateCmd.Flags().String("ms-sql", "", "Microsoft SQL Server")
}
