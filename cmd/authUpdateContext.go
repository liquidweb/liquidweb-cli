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

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/types/cmd"
	"github.com/liquidweb/liquidweb-cli/validate"
)

var authUpdateContextCmd = &cobra.Command{
	Use:   "update-context",
	Short: "Update an existing auth context",
	Long: `Update an existing auth context.

If you've never setup any contexts, check "auth init".`,
	Run: func(cmd *cobra.Command, args []string) {
		contextName, _ := cmd.Flags().GetString("context")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		url, _ := cmd.Flags().GetString("url")
		timeout, _ := cmd.Flags().GetInt("timeout")
		setInsecure, _ := cmd.Flags().GetBool("set-insecure")
		setSecure, _ := cmd.Flags().GetBool("set-secure")

		if username == "" && password == "" && url == "" && timeout == -1 &&
			!setInsecure && !setSecure {
			lwCliInst.Die(fmt.Errorf("must pass something to update"))
		}

		if setInsecure && setSecure {
			lwCliInst.Die(fmt.Errorf("cant set insecure and secure"))
		}

		contexts := lwCliInst.Viper.GetStringMap("liquidweb.api.contexts")
		if _, exists := contexts[contextName]; !exists {
			lwCliInst.Die(fmt.Errorf("context with name [%s] doesnt exist", contextName))
		}

		var authContext cmdTypes.AuthContext
		if err := instance.CastFieldTypes(contexts[contextName], &authContext); err != nil {
			lwCliInst.Die(err)
		}

		validateFields := map[interface{}]interface{}{}

		if username != "" {
			authContext.Username = username
			validateFields[username] = "NonEmptyString"
		}
		if password != "" {
			authContext.Password = password
			validateFields[password] = "NonEmptyString"
		}
		if url != "" {
			authContext.Url = url
			validateFields[url] = "HttpsLiquidwebUrl"
		}
		if timeout != -1 {
			authContext.Timeout = timeout
			validateFields[timeout] = "PositiveInt"
		}
		if setSecure {
			authContext.Insecure = false
		}
		if setInsecure {
			authContext.Insecure = true
		}

		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s", contextName), map[string]interface{}{
			"contextname": contextName,
			"username":    authContext.Username,
			"password":    authContext.Password,
			"url":         authContext.Url,
			"insecure":    authContext.Insecure,
			"timeout":     authContext.Timeout,
		})

		if err := lwCliInst.Viper.WriteConfig(); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Printf("Updated context [%s]\n", contextName)
	},
}

func init() {
	authCmd.AddCommand(authUpdateContextCmd)

	authUpdateContextCmd.Flags().String("context", "", "name of context")
	authUpdateContextCmd.Flags().String("username", "", "api username")
	authUpdateContextCmd.Flags().String("password", "", "password for username")
	authUpdateContextCmd.Flags().String("url", "", "api url")
	authUpdateContextCmd.Flags().Int("timeout", -1, "api timeout value")
	authUpdateContextCmd.Flags().Bool("set-insecure", false, "enable insecure SSL validation of api url")
	authUpdateContextCmd.Flags().Bool("set-secure", false, "enable secure SSL validation of api url")

	if err := authUpdateContextCmd.MarkFlagRequired("context"); err != nil {
		lwCliInst.Die(err)
	}
}
