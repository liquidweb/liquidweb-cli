/*
Copyright © 2019 LiquidWeb

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
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/utils"
)

var authAddContextCmd = &cobra.Command{
	Use:   "add-context",
	Short: "Add a context to an existing configuration",
	Long: `Add a context to an existing configuration.

Use this if you've already setup contexts with "auth init".`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		contextName, _ := cmd.Flags().GetString("context")
		url, _ := cmd.Flags().GetString("api-url")
		insecure, _ := cmd.Flags().GetBool("insecure")
		timeout, _ := cmd.Flags().GetInt("timeout")

		contextName = strings.ToLower(contextName)

		file, err := getExpectedConfigPath()
		if err != nil {
			lwCliInst.Die(err)
		}
		if !utils.FileExists(file) {
			f, err := os.Create(filepath.Clean(file))
			if err != nil {
				lwCliInst.Die(err)
			}
			if err := f.Close(); err != nil {
				lwCliInst.Die(err)
			}
			if err := os.Chmod(file, 0600); err != nil {
				lwCliInst.Die(err)
			}
		}

		contexts := lwCliInst.Viper.GetStringMap("liquidweb.api.contexts")
		if _, exists := contexts[contextName]; exists {
			lwCliInst.Die(fmt.Errorf("context with name [%s] already exists", contextName))
		}

		lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s", contextName), map[string]interface{}{
			"contextname": contextName,
			"username":    username,
			"password":    password,
			"url":         url,
			"insecure":    insecure,
			"timeout":     timeout,
		})

		if err := lwCliInst.Viper.WriteConfig(); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Created context [%s]\n", contextName)
	},
}

func init() {
	authCmd.AddCommand(authAddContextCmd)

	authAddContextCmd.Flags().String("context", "", "name for new context")
	authAddContextCmd.Flags().String("username", "", "username to authenticate with")
	authAddContextCmd.Flags().String("password", "", "password for given username")
	authAddContextCmd.Flags().Bool("insecure", false, "whether or not to perform SSL validation on api url")
	authAddContextCmd.Flags().String("api-url", "https://api.liquidweb.com", "API URL to use")
	authAddContextCmd.Flags().Int("timeout", 30, "timeout value when communicating with api-url")

	if err := authAddContextCmd.MarkFlagRequired("username"); err != nil {
		lwCliInst.Die(err)
	}
	if err := authAddContextCmd.MarkFlagRequired("password"); err != nil {
		lwCliInst.Die(err)
	}
}
