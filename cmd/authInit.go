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
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/types/cmd"
	"github.com/liquidweb/liquidweb-cli/types/errors"
)

var authInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Specify your LiquidWeb API credentials from a blank slate",
	Long: `Specify your LiquidWeb API credentials from a blank slate.

Intended to be ran for initial setup only.`,
	Run: func(cmd *cobra.Command, args []string) {
		writeEmptyConfig()

		if err := setAuthDataInteractively(); err != nil {
			lwCliInst.Die(err)
		}
	},
}

func init() {
	authCmd.AddCommand(authInitCmd)
}

func setAuthDataInteractively() error {
	_, err := fetchAuthDataInteractively()
	if err != nil {
		return err
	}

	if err := lwCliInst.Viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func fetchAuthDataInteractively() ([]cmdTypes.AuthContext, error) {
	var contexts []cmdTypes.AuthContext

	if !terminal.IsTerminal(0) || !terminal.IsTerminal(1) {
		return contexts, errorTypes.UnknownTerminal
	}
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return contexts, err
	}
	defer terminal.Restore(0, oldState)
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	term := terminal.NewTerminal(screen, "")
	term.SetPrompt(cast.ToString(term.Escape.Blue) + " > " + cast.ToString(term.Escape.Reset))

	moreAdds := true
	for moreAdds {
		// name context
		fmt.Printf("Name this context: ")
		contextNameBytes, err := term.ReadLine()
		if err != nil {
			return contexts, err
		}
		contextNameString := cast.ToString(contextNameBytes)
		for contextNameString == "" {
			fmt.Printf("context name cannot be blank. Please name this context: ")
			contextNameBytes, err = term.ReadLine()
			if err != nil {
				return contexts, err
			}
			contextNameString = cast.ToString(contextNameBytes)
		}
		lwCliInst.Viper.Set(
			fmt.Sprintf("liquidweb.api.contexts.%s.contextname",
				contextNameString), contextNameString)

		// username
		fmt.Print("LiquidWeb username: ")
		usernameBytes, err := term.ReadLine()
		if err != nil {
			return contexts, err
		}
		usernameString := cast.ToString(usernameBytes)
		lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s.username",
			contextNameString), usernameString)

		// password
		passwordBytes, err := term.ReadPassword("LiquidWeb password: ")
		if err != nil {
			return contexts, err
		}
		passwordString := cast.ToString(passwordBytes)
		lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s.password",
			contextNameString), passwordString)

		// url
		fmt.Printf("API URL (hit enter for default): ")
		urlBytes, err := term.ReadLine()
		if err != nil {
			return contexts, err
		}
		urlString := cast.ToString(urlBytes)
		if urlString == "" {
			urlString = "https://api.liquidweb.com"
		}
		lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s.url",
			contextNameString), urlString)

		// insecure ssl validation
		fmt.Printf("Insecure SSL Validation (yes/no) (hit enter for default): ")
		insecureBytes, err := term.ReadLine()
		if err != nil {
			return contexts, err
		}
		insecureString := cast.ToString(insecureBytes)
		insecure := false
		if insecureString == "yes" {
			insecure = true
		}
		lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s.insecure",
			contextNameString), insecure)

		// timeout
		fmt.Printf("API timeout (hit enter for default): ")
		timeoutBytes, err := term.ReadLine()
		if err != nil {
			return contexts, err
		}
		timeoutInt := cast.ToInt(timeoutBytes)
		if timeoutInt == 0 {
			timeoutInt = 30
		}
		lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s.timeout",
			contextNameString), timeoutInt)

		// make current context?
		fmt.Printf("Make current context? (yes/no)")
		makeCurrentContextBytes, err := term.ReadLine()
		if err != nil {
			return contexts, err
		}
		makeCurrentContextString := cast.ToString(makeCurrentContextBytes)
		if makeCurrentContextString == "yes" {
			lwCliInst.Viper.Set("liquidweb.api.current_context", contextNameString)
		}

		// more contexts to add ?
		fmt.Printf("Add another context? (yes/no): ")
		moreContextsBytes, err := term.ReadLine()
		if err != nil {
			return contexts, err
		}
		if cast.ToString(moreContextsBytes) == "no" {
			moreAdds = false
		}

		// save entry data
		context := cmdTypes.AuthContext{
			ContextName: contextNameString,
			Username:    usernameString,
			Password:    passwordString,
			Url:         urlString,
			Insecure:    insecure,
			Timeout:     timeoutInt,
		}

		contexts = append(contexts, context)
	}

	return contexts, err
}

func writeEmptyConfig() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	cfgFile := fmt.Sprintf("%s/.liquidweb-cli.yaml", homedir)
	fmt.Printf("using config file %s\n", cfgFile)
	f, err := os.Create(cfgFile)
	if err != nil {
		return err
	}
	if err := os.Chmod(cfgFile, 0600); err != nil {
		return err
	}

	f.Close()

	return nil
}
