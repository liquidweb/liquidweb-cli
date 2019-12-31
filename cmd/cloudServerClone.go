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
)

var cloudServerCloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a Cloud Server",
	Long: `Clone a Cloud Server.

Returns the information about the newly created clone.

All of the optional fields are defaulted to the values on the original server if
they aren't received. For cloning to a private parent, include the uniq_id of the
parent server to be cloned to, along with the memory/diskspace/vcpu amounts (if
different from the original).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("clone called")
	},
}

func init() {
	cloudServerCmd.AddCommand(cloudServerCloneCmd)
}
