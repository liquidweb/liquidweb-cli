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
	"os"

	"github.com/spf13/cobra"
)

var cloudServerBlockStorageOptimizedCmd = &cobra.Command{
	Use:   "block-storage-optimized",
	Short: "Cloud Block Storage Optimized actions",
	Long: `Cloud Block Storage Optimized actions.

A Cloud Server uses more memory when using Cloud Block Storage volumes. Normal
Cloud Servers have enough pad for this to not be an issue, but Cloud Dedicated
purposely runs very tight. To work around this we can reduce the RAM allocated
to the Cloud Server to give more to the hypervisor. We call this Cloud Block
Storage Optimized.

Enabling or disabling Cloud Block Storage will cause your Cloud Server to reboot.

For a full list of capabilities, please refer to the "Available Commands" section.`,

	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			lwCliInst.Die(err)
		}
		os.Exit(1)
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerBlockStorageOptimizedCmd)
}
