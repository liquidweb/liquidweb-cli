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

	"github.com/liquidweb/liquidweb-cli/instance"
)

var cloudTemplateRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a Cloud Template on a Cloud Server",
	Long:  `Restore a Cloud Template on a Cloud Server.`,
	Run: func(cmd *cobra.Command, args []string) {
		params := &instance.CloudTemplateRestoreParams{}

		params.UniqId, _ = cmd.Flags().GetString("uniq-id")
		params.Template, _ = cmd.Flags().GetString("template")

		result, err := lwCliInst.CloudTemplateRestore(params)
		if err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Restoring template! %s\n", result)
		fmt.Printf("\tcheck progress with 'cloud server status --uniq-id %s'\n", params.UniqId)
	},
}

func init() {
	cloudTemplateCmd.AddCommand(cloudTemplateRestoreCmd)

	cloudTemplateRestoreCmd.Flags().String("uniq-id", "", "uniq-id of Cloud Server")
	cloudTemplateRestoreCmd.Flags().String("template", "", "name of template to restore")

	cloudTemplateRestoreCmd.MarkFlagRequired("uniq-id")
	cloudTemplateRestoreCmd.MarkFlagRequired("template")
}
