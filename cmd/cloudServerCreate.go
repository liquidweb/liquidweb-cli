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
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/instance"
)

var cloudServerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Cloud Server",
	Long: `Create a Cloud Server.

Requires various flags. Please see the flag section of help.`,
	Run: func(cmd *cobra.Command, args []string) {
		verboseFlag, _ := cmd.Flags().GetBool("verbose")
		jsonFlag, _ := cmd.Flags().GetBool("json")
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
		antivirusFlag, _ := cmd.Flags().GetString("antivirus")
		msSqlFlag, _ := cmd.Flags().GetString("ms_sql")
		privateParentFlag, _ := cmd.Flags().GetString("private-parent")
		passwordFlag, _ := cmd.Flags().GetString("password")
		memoryFlag, _ := cmd.Flags().GetInt("memory")
		diskspaceFlag, _ := cmd.Flags().GetInt("diskspace")
		vcpuFlag, _ := cmd.Flags().GetInt("vcpu")

		// sanity check flags
		if zoneFlag == 0 {
			lwCliInst.Die(fmt.Errorf("--zone is a required flag"))
		}
		if configIdFlag == 0 && privateParentFlag == "" {
			lwCliInst.Die(fmt.Errorf("--config_id is a required flag without --private-parent"))
		}
		if templateFlag == "" {
			// TODO: not required when creating from a backup or an image
			lwCliInst.Die(fmt.Errorf("--template is a required flag"))
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
				"Template":  templateFlag,
				"ExtraIp": map[string]interface{}{
					"value": ipsFlag,
					"count": 0,
				},
				"LiquidWebBackupPlan": backupPlanFlag,
			},
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

		if antivirusFlag != "" {
			createArgs["antivirus"] = antivirusFlag
		}
		if msSqlFlag != "" {
			createArgs["ms_sql"] = msSqlFlag
		}

		if backupPlanFlag == "Quota" {
			createArgs["features"].(map[string]interface{})["BackupQuota"] = backupPlanQuotaFlag
		}

		if publicSshKeyContents != "" {
			createArgs["public_ssh_key"] = publicSshKeyContents
		}

		if verboseFlag {
			pr, err := lwCliInst.JsonEncodeAndPrettyPrint(createArgs)
			if err == nil {
				fmt.Println("createArgs:")
				fmt.Println(pr)
			}
		}

		result, err := lwCliInst.LwApiClient.Call("bleed/server/create", createArgs)
		if err != nil {
			lwCliInst.Die(err)
		}

		if jsonFlag {
			pretty, err := lwCliInst.JsonEncodeAndPrettyPrint(result)
			if err != nil {
				lwCliInst.Die(err)
			}
			fmt.Printf(pretty)
			os.Exit(0)
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

	cloudServerCreateCmd.Flags().Bool("verbose", false, "provide verbose output")
	cloudServerCreateCmd.Flags().Bool("json", false, "output in json format")
	cloudServerCreateCmd.Flags().String("template", "", "name of the template to use")
	cloudServerCreateCmd.Flags().String("type", "SS.VPS", "some examples of types; SS.VPS, SS.VPS.WIN, SS.VM, SS.VM.WIN")
	cloudServerCreateCmd.Flags().String("hostname", randomHostname, "hostname to set")
	cloudServerCreateCmd.Flags().Int("ips", 1, "amount of IP addresses")
	// TODO pool_ips
	cloudServerCreateCmd.Flags().String("public-ssh-key", sshPubKeyFile,
		"path to file containing the public ssh key you wish to be on the new Cloud Server")
	cloudServerCreateCmd.Flags().Int("config_id", 0, "config_id to use")
	cloudServerCreateCmd.Flags().String("backup-plan", "None", "Cloud Server backup plan to use")
	cloudServerCreateCmd.Flags().Int("backup-plan-quota", 300, "Quota amount. Should only be used with '--backup-plan Quota'")
	cloudServerCreateCmd.Flags().String("bandwidth", "SS.10000", "bandwidth package to use")
	cloudServerCreateCmd.Flags().Int("zone", 0, "Cloud server zone to create in")
	cloudServerCreateCmd.Flags().String("password", randomPassword, "root or administrator password to set")

	// private parent specific
	cloudServerCreateCmd.Flags().String("private-parent", "",
		"name or uniq_id of the private-parent. Must use when creating a Cloud Server on a private parent.")
	cloudServerCreateCmd.Flags().Int("memory", -1, "memory (ram) value use with --private-parent")
	cloudServerCreateCmd.Flags().Int("diskspace", -1, "diskspace value use with --private-parent")
	cloudServerCreateCmd.Flags().Int("vcpu", -1, "vcpu value use with --private-parent")

	// windows specific
	cloudServerCreateCmd.Flags().String("antivirus", "", "antivirus string")
	cloudServerCreateCmd.Flags().String("ms-sql", "", "ms_sql string")
}
