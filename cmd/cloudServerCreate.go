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
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/instance"
)

var cloudServerCreateCmdPoolIpsFlag []string
var cloudServerCreateCmdPool6IpsFlag []string

var cloudServerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Cloud Server",
	Long: `Create a Cloud Server.

Requires various flags. Please see the flag section of help.

Examples:

# Create a Cloud Server on a Private Parent named "private"
'cloud server create --private-parent private --memory 1024 --diskspace 40 --vcpu 2 --template DEBIAN_10_UNMANAGED'

# Create a Cloud Server on config-id 1
'cloud server create --config-id 1 --template DEBIAN_10_UNMANAGED --zone 40460'

# Create a Cloud Server from image id 111
'cloud server create --image-id 111 --zone 40460 --config-id 1'

# Create a Cloud Server from backup id 111
'cloud server create --backup-id 111 --zone 40460 --config-id 1'

These examples use default values for various flags, such as password, type, ssh-key, hostname, etc.

For a list of Templates, Configs, and Region/Zones, see 'cloud server options --configs --templates --zones'
For a list of images, see 'cloud images list'
For a list of backups, see 'cloud backups list'

Plan Example:

---
cloud:
   server:
      create:
         - type: "SS.VPS.WIN"
           password: "1fk4ds$jktl43u90dsa"
           template: "WINDOWS_2019_UNMANAGED"
           zone: 40460
           hostname: "db1.dev.addictmud.org"
           ips: 1
           ip6s: 1
           public-ssh-key: ""
           config-id: 88
           backup-days: 5
           bandwidth: "SS.5000"
           backup-id: -1
           image-id: -1
           pool-ips:
              - "10.111.12.13"
              - "10.12.13.14"
           pool6-ips:
              - "2607:fad0:3714:350c::/64"
           private-parent: "my pp"
           memory: 0
           diskspace: 0
           vcpu: 0
           winav: ""
           ms-sql: ""

lw plan --file /tmp/cloud.server.create.yaml
`,
	Run: func(cmd *cobra.Command, args []string) {
		params := &instance.CloudServerCreateParams{}

		params.Template, _ = cmd.Flags().GetString("template")
		params.Type, _ = cmd.Flags().GetString("type")
		params.Hostname, _ = cmd.Flags().GetString("hostname")
		params.Ips, _ = cmd.Flags().GetInt("ips")
		params.Ip6s, _ = cmd.Flags().GetInt("ip6s")
		pubSshKey, _ := cmd.Flags().GetString("public-ssh-key")
		params.ConfigId, _ = cmd.Flags().GetInt("config-id")
		params.BackupDays, _ = cmd.Flags().GetInt("backup-days")
		params.BackupQuota, _ = cmd.Flags().GetInt("backup-quota")
		params.Bandwidth, _ = cmd.Flags().GetString("bandwidth")
		params.Zone, _ = cmd.Flags().GetInt64("zone")
		params.WinAv, _ = cmd.Flags().GetString("winav")
		params.MsSql, _ = cmd.Flags().GetString("ms-sql")
		params.PrivateParent, _ = cmd.Flags().GetString("private-parent")
		params.Password, _ = cmd.Flags().GetString("password")
		params.Memory, _ = cmd.Flags().GetInt("memory")
		params.Diskspace, _ = cmd.Flags().GetInt("diskspace")
		params.Vcpu, _ = cmd.Flags().GetInt("vcpu")
		params.BackupId, _ = cmd.Flags().GetInt("backup-id")
		params.ImageId, _ = cmd.Flags().GetInt("image-id")
		params.PoolIps = cloudServerCreateCmdPoolIpsFlag
		params.Pool6Ips = cloudServerCreateCmdPool6IpsFlag

		sshPkeyContents, err := ioutil.ReadFile(filepath.Clean(pubSshKey))
		if err == nil {
			params.PublicSshKey = cast.ToString(sshPkeyContents)
		}

		uniqId, err := lwCliInst.CloudServerCreate(params)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf(
			"Cloud server with uniq-id [%s] creating. Check status with 'cloud server status --uniq-id %s'\n",
			uniqId, uniqId)
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerCreateCmd)

	var sshPubKeyFile string
	home, err := homedir.Dir()
	if err == nil {
		sshPubKeyFile = fmt.Sprintf("%s/.ssh/id_rsa.pub", home)
	}

	cloudServerCreateCmd.Flags().String("template", "", "template to use (see 'cloud server options --templates')")
	cloudServerCreateCmd.Flags().String("type", "SS.VPS", "some examples of types; SS.VPS, SS.VPS.WIN, SS.VM, SS.VM.WIN")
	cloudServerCreateCmd.Flags().String("hostname", "", "hostname to set")
	cloudServerCreateCmd.Flags().Int("ips", 1, "amount of IPv4 addresses")
	cloudServerCreateCmd.Flags().Int("ip6s", 0, "amount of IPv6 /64s")
	cloudServerCreateCmd.Flags().String("public-ssh-key", sshPubKeyFile,
		"path to file containing the public ssh key you wish to be on the new Cloud Server")
	cloudServerCreateCmd.Flags().Int("config-id", 0, "config-id to use")
	cloudServerCreateCmd.Flags().Int("backup-days", -1, "Enable daily backup plan. This is the amount of days to keep a backup")
	cloudServerCreateCmd.Flags().Int("backup-quota", -1, "Enable quota backup plan. This is the total amount of GB to keep.")
	cloudServerCreateCmd.Flags().String("bandwidth", "SS.10000", "bandwidth package to use")
	cloudServerCreateCmd.Flags().Int64("zone", 0, "zone (id) to create new Cloud Server in (see 'cloud server options --zones')")
	cloudServerCreateCmd.Flags().String("password", "", "root or administrator password to set")

	cloudServerCreateCmd.Flags().Int("backup-id", -1, "id of cloud backup to create from (see 'cloud backup list')")
	cloudServerCreateCmd.Flags().Int("image-id", -1, "id of cloud image to create from (see 'cloud image list')")

	cloudServerCreateCmd.Flags().StringSliceVar(&cloudServerCreateCmdPoolIpsFlag, "pool-ips", []string{},
		"IPv4 ips from your IP Pool separated by ',' to assign to the new Cloud Server")

	cloudServerCreateCmd.Flags().StringSliceVar(&cloudServerCreateCmdPool6IpsFlag, "pool6-ips", []string{},
		"IPv6 assignments from your IP Pool separated by ',' to assign to the new Cloud Server")

	// private parent specific
	cloudServerCreateCmd.Flags().String("private-parent", "",
		"name or uniq-id of the private-parent. Must use when creating a Cloud Server on a private parent.")
	cloudServerCreateCmd.Flags().Int("memory", -1, "memory (ram) value use with --private-parent")
	cloudServerCreateCmd.Flags().Int("diskspace", -1, "diskspace value use with --private-parent")
	cloudServerCreateCmd.Flags().Int("vcpu", -1, "vcpu value use with --private-parent")

	// windows specific
	cloudServerCreateCmd.Flags().String("winav", "", "Use only with Windows Servers. Typically (None or NOD32) for value when set")
	cloudServerCreateCmd.Flags().String("ms-sql", "", "Microsoft SQL Server")
}
