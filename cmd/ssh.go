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

Gives you an interactive shell to your Server.`,
	Run: func(cmd *cobra.Command, args []string) {
		params := &instance.SshParams{}

		params.Host, _ = cmd.Flags().GetString("host")
		params.PrivateKey, _ = cmd.Flags().GetString("private-ssh-key")
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
	sshCmd.Flags().String("private-ssh-key", "", "path to a specific/non default ssh private key to use")
	sshCmd.Flags().Bool("agent-forwarding", false, "whether or not to enable ssh agent forwarding")
	sshCmd.Flags().String("user", "root", "username to use for the ssh connection")
	sshCmd.Flags().String("command", "", "run this command and exit rather than start an interactive shell")

	if err := sshCmd.MarkFlagRequired("host"); err != nil {
		lwCliInst.Die(err)
	}
}
