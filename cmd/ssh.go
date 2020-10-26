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
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/instance"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH to a Server",
	Long: `SSH to a Server.

Starts an interactive SSH session to your server. If --command is passed, a non interactive
session is started running only your passed command.

Examples:

* SSH by hostname accepting all defaults:
- lw ssh --host dexg.ulxy5e656r.io

* SSH by uniq-id accepting all defaults:
- lw ssh --host ABC123

* SSH by uniq-id making use of all flags:
- lw ssh --host ABC123 --agent-forwarding --port 2222 --private-key-file /home/myself/.ssh-alt/id_rsa \
    --user amanda --command "ps faux && free -m"

Plan Examples:

---
ssh:
  - host: nd00.ltv1wv76kc.io
    command: "free -m"
    user: "root"
    private-key-file: "/home/myself/.ssh/id_rsa"
    agent-forwarding: true
    port: 22
  - host: PPB4NZ
    command: "hostname && free -m"

lw plan --file /tmp/ssh.yaml
`,
	Run: func(cmd *cobra.Command, args []string) {
		params := &instance.SshParams{}

		params.Host, _ = cmd.Flags().GetString("host")
		params.PrivateKeyFile, _ = cmd.Flags().GetString("private-key-file")
		params.User, _ = cmd.Flags().GetString("user")
		params.AgentForwarding, _ = cmd.Flags().GetBool("agent-forwarding")
		params.Port, _ = cmd.Flags().GetInt("port")
		params.Command, _ = cmd.Flags().GetString("command")

		if err := lwCliInst.Ssh(params); err != nil {
			lwCliInst.Die(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)

	sshCmd.Flags().String("host", "", "uniq-id or hostname for the Server")
	sshCmd.Flags().Int("port", 22, "ssh port to use")
	sshCmd.Flags().String("private-key-file", "", "path to a specific/non default ssh private key to use")
	sshCmd.Flags().Bool("agent-forwarding", false, "whether or not to enable ssh agent forwarding")
	sshCmd.Flags().String("user", "root", "username to use for the ssh connection")
	sshCmd.Flags().String("command", "", "run this command and exit rather than start an interactive shell")

	if err := sshCmd.MarkFlagRequired("host"); err != nil {
		lwCliInst.Die(err)
	}
}
