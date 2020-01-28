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
	"bytes"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/types/cmd"
	"github.com/liquidweb/liquidweb-cli/types/errors"
	"github.com/liquidweb/liquidweb-cli/utils"
)

var authInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Specify your LiquidWeb API credentials from a blank slate",
	Long: `Specify your LiquidWeb API credentials from a blank slate.

Intended to be ran for initial setup only.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := setAuthDataInteractively(); err != nil {
			lwCliInst.Die(err)
		}
	},
}

func init() {
	authCmd.AddCommand(authInitCmd)
}

func setAuthDataInteractively() error {
	_, writeConfig, err := fetchAuthDataInteractively()
	if err != nil {
		return err
	}

	if writeConfig {
		if err := lwCliInst.Viper.WriteConfig(); err != nil {
			return err
		}
	}

	return nil
}

func fetchAuthDataInteractively() ([]cmdTypes.AuthContext, bool, error) {
	var contexts []cmdTypes.AuthContext

	if !terminal.IsTerminal(0) || !terminal.IsTerminal(1) {
		return contexts, false, errorTypes.UnknownTerminal
	}
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return contexts, false, err
	}
	defer terminal.Restore(0, oldState)
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	term := terminal.NewTerminal(screen, "")
	term.SetPrompt(cast.ToString(term.Escape.Blue) + " > " + cast.ToString(term.Escape.Reset))

	moreAdds := false

	// warn before deleting config
	var haveProceedAnswer bool
	for !haveProceedAnswer {
		term.Write([]byte("Warning: This will delete all auth contexts. Continue (yes/[no])?: "))
		proceedBytes, err := term.ReadLine()
		if err != nil {
			return contexts, false, err
		}
		proceedString := cast.ToString(proceedBytes)
		if proceedString != "yes" && proceedString != "no" && proceedString != "" {
			term.Write([]byte("invalid input.\n"))
			continue
		}

		if proceedString == "yes" {
			moreAdds = true
			haveProceedAnswer = true
		} else if proceedString == "" || proceedString == "no" {
			haveProceedAnswer = true
			moreAdds = false
		}
	}

	// return if user didnt acknowledge to proceed
	if !moreAdds {
		return contexts, false, nil
	}

	// if user consented to proceed, clear config
	writeEmptyConfig()
	cfgFile, err := getExpectedConfigPath()
	if err != nil {
		return contexts, false, err
	}
	if utils.FileExists(cfgFile) {
		if err := os.Remove(cfgFile); err != nil {
			lwCliInst.Die(err)
		}
		f, err := os.Create(cfgFile)
		if err != nil {
			lwCliInst.Die(err)
		}
		f.Close()
		if err := os.Chmod(cfgFile, 0600); err != nil {
			lwCliInst.Die(err)
		}

		lwCliInst.Viper.ReadConfig(bytes.NewBuffer([]byte{}))
	}

	for moreAdds {
		var (
			contextNameAnswer            string
			haveContextNameAnswer        bool
			usernameAnswer               string
			haveUsernameAnswer           bool
			passwordAnswer               string
			havePasswordAnswer           bool
			haveMakeCurrentContextAnswer bool
			haveMoreContextsToAddAnswer  bool
		)

		// context name
		for !haveContextNameAnswer {
			term.Write([]byte("Name this context: "))
			contextNameBytes, err := term.ReadLine()
			if err != nil {
				return contexts, false, err
			}
			contextNameAnswer = cast.ToString(contextNameBytes)
			if contextNameAnswer == "" {
				term.Write([]byte("context name cannot be blank.\n"))
			} else {
				haveContextNameAnswer = true
				lwCliInst.Viper.Set(fmt.Sprintf(
					"liquidweb.api.contexts.%s.contextname", contextNameAnswer), contextNameAnswer)
			}
		}

		// username
		for !haveUsernameAnswer {
			term.Write([]byte("LiquidWeb username: "))
			usernameBytes, err := term.ReadLine()
			if err != nil {
				return contexts, false, err
			}
			usernameAnswer = cast.ToString(usernameBytes)
			if usernameAnswer == "" {
				term.Write([]byte("username cannot be blank.\n"))
			} else {
				haveUsernameAnswer = true
				lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s.username",
					contextNameAnswer), usernameAnswer)
			}
		}

		// password
		for !havePasswordAnswer {
			passwordBytes, err := term.ReadPassword("LiquidWeb password: ")
			if err != nil {
				return contexts, false, err
			}
			passwordAnswer = cast.ToString(passwordBytes)
			if passwordAnswer == "" {
				term.Write([]byte("password cannot be blank.\n"))
			} else {
				havePasswordAnswer = true
				lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s.password",
					contextNameAnswer), passwordAnswer)
			}
		}

		// make current context?
		for !haveMakeCurrentContextAnswer {
			term.Write([]byte("Make current context? ([yes]/no)"))
			makeCurrentContextBytes, err := term.ReadLine()
			if err != nil {
				return contexts, false, err
			}
			makeCurrentContextString := cast.ToString(makeCurrentContextBytes)
			if makeCurrentContextString != "" && makeCurrentContextString != "yes" && makeCurrentContextString != "no" {
				term.Write([]byte("invalid input.\n"))
				continue
			}

			haveMakeCurrentContextAnswer = true
			if makeCurrentContextString == "yes" || makeCurrentContextString == "" {
				lwCliInst.Viper.Set("liquidweb.api.current_context", contextNameAnswer)
			}
		}

		// more contexts to add ?
		for !haveMoreContextsToAddAnswer {
			term.Write([]byte("Add another context? (yes/[no]): "))
			moreContextsBytes, err := term.ReadLine()
			if err != nil {
				return contexts, false, err
			}

			answer := cast.ToString(moreContextsBytes)
			if answer != "" && answer != "yes" && answer != "no" {
				term.Write([]byte("invalid input.\n"))
				continue
			}

			if answer == "no" || answer == "" {
				moreAdds = false
				haveMoreContextsToAddAnswer = true
			}

			haveMoreContextsToAddAnswer = true
		}

		// if you can't use these defaults, see `auth update-context` to change it later
		defaultUrl := "https://api.liquidweb.com"
		lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s.url",
			contextNameAnswer), defaultUrl)
		defaultTimeout := 90
		lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s.timeout",
			contextNameAnswer), defaultTimeout)
		defaultInsecure := false
		lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s.insecure",
			contextNameAnswer), defaultInsecure)

		// return entry data
		context := cmdTypes.AuthContext{
			ContextName: contextNameAnswer,
			Username:    usernameAnswer,
			Password:    passwordAnswer,
			Url:         defaultUrl,
			Insecure:    defaultInsecure,
			Timeout:     defaultTimeout,
		}

		contexts = append(contexts, context)
	}

	return contexts, true, err
}

func getExpectedConfigPath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cfgFile := fmt.Sprintf("%s/.liquidweb-cli.yaml", homedir)

	return cfgFile, nil
}

func writeEmptyConfig() error {
	cfgFile, err := getExpectedConfigPath()
	if err != nil {
		return err
	}

	f, err := os.Create(cfgFile)
	if err != nil {
		return err
	}
	f.Close()

	if err := os.Chmod(cfgFile, 0600); err != nil {
		return err
	}

	return nil
}
