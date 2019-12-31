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

	"github.com/spf13/cobra"
)

var authAddContextCmd = &cobra.Command{
	Use:   "add-context",
	Short: "Add a context to an existing configuration",
	Long: `Add a context to an existing configuration.

Use this if you've already setup contexts with "auth init".`,
	Run: func(cmd *cobra.Command, args []string) {
		requiredStrFlags := []string{
			"context-name",
			"username",
			"password",
		}
		for _, requiredStrFlag := range requiredStrFlags {
			val, err := cmd.Flags().GetString(requiredStrFlag)
			if err != nil {
				lwCliInst.Die(err)
			}
			if val == "" {
				lwCliInst.Die(fmt.Errorf("required flag [%s] was not provided", requiredStrFlag))
			}

		}
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			lwCliInst.Die(err)
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			lwCliInst.Die(err)
		}
		contextName, err := cmd.Flags().GetString("context-name")
		if err != nil {
			lwCliInst.Die(err)
		}
		url, err := cmd.Flags().GetString("api-url")
		if err != nil {
			lwCliInst.Die(err)
		}
		insecure, err := cmd.Flags().GetBool("insecure")
		if err != nil {
			lwCliInst.Die(err)
		}
		timeout, err := cmd.Flags().GetInt("timeout")
		if err != nil {
			lwCliInst.Die(err)
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

	authAddContextCmd.Flags().String("context-name", "", "name of context")
	authAddContextCmd.Flags().String("username", "", "username to authenticate with")
	authAddContextCmd.Flags().String("password", "", "password for given username")
	authAddContextCmd.Flags().Bool("insecure", false, "whether or not to perform SSL validation on api url")
	authAddContextCmd.Flags().String("api-url", "https://api.liquidweb.com", "API URL to use")
	authAddContextCmd.Flags().Int("timeout", 30, "timeout value when communicating with api-url")
}
