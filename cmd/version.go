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

// default build-time variables; these are overridden via ldflags
var (
	Version   = "unknown-version"
	GitCommit = "unknown-commit"
	BuildTime = "unknown-buildtime"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show build information",
	Long: `Show build information.

This information should be provided with any bug report.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("LiquidWeb CLI Build Details\n\n")
		fmt.Printf("  Build Time: %s\n", BuildTime)
		fmt.Printf("  Version: %s\n", Version)
		fmt.Printf("  Git commit: %s\n\n", GitCommit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
