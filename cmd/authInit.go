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
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"

	cmdTypes "github.com/liquidweb/liquidweb-cli/types/cmd"
	"github.com/liquidweb/liquidweb-cli/utils"
	"github.com/liquidweb/liquidweb-cli/validate"
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
	writeConfig, err := fetchAuthDataInteractively()
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

func fetchAuthDataInteractively() (writeConfig bool, err error) {
	var (
		moreAdds          bool
		haveProceedAnswer bool
	)

	for !haveProceedAnswer {
		f := func(d prompt.Document) []prompt.Suggest {
			s := []prompt.Suggest{
				{Text: "yes", Description: "delete all auth contexts"},
				{Text: "no", Description: "keep my auth contexts and exit"},
			}
			return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
		}
		fmt.Println("Warning: If you proceed, any set auth context will be deleted.")
		fmt.Print("\nAre you sure? ")
		answer := strings.ToLower(prompt.Input("> ", f, prompt.OptionShowCompletionAtStart()))
		if answer == "yes" {
			moreAdds = true
			haveProceedAnswer = true
		} else if answer == "no" {
			haveProceedAnswer = true
			moreAdds = false
		}
	}

	if !moreAdds {
		return
	}

	userInputComplete := make(chan bool)
	userInputError := make(chan error)
	userInputExitEarly := make(chan bool)
	userInputContext := make(chan cmdTypes.AuthContext)

	fmt.Println("\nTo exit early, type 'exit'\n")

	go func() {
	WHILEMOREADDS:
		for moreAdds {
			var (
				context                      cmdTypes.AuthContext
				haveContextNameAnswer        bool
				haveUsernameAnswer           bool
				havePasswordAnswer           bool
				haveMakeCurrentContextAnswer bool
				haveMoreContextsToAddAnswer  bool
			)

			// context name
			for !haveContextNameAnswer {
				fmt.Print("Name this context: ")
				answer := strings.ToLower(prompt.Input("> ", func(d prompt.Document) (s []prompt.Suggest) { return }))
				if answer == "exit" {
					userInputExitEarly <- true
					break WHILEMOREADDS
				} else if answer != "" {
					haveContextNameAnswer = true
					context.ContextName = answer
				}
			}

			// username
			for !haveUsernameAnswer {
				fmt.Print("LiquidWeb username: ")
				answer := prompt.Input("> ", func(d prompt.Document) (s []prompt.Suggest) { return })
				if answer == "exit" {
					userInputExitEarly <- true
					break WHILEMOREADDS
				} else if answer != "" {
					haveUsernameAnswer = true
					context.Username = answer
				}
			}

			// password
			for !havePasswordAnswer {
				fmt.Print("LiquidWeb password: ")
				passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
				if err != nil {
					userInputError <- err
					break WHILEMOREADDS
				}
				answer := string(passwordBytes)
				if answer == "exit" {
					userInputExitEarly <- true
					break WHILEMOREADDS
				} else if answer != "" {
					havePasswordAnswer = true
					context.Password = answer
				}
				fmt.Println("")
			}

			// make current context?
			for !haveMakeCurrentContextAnswer {
				f := func(d prompt.Document) []prompt.Suggest {
					s := []prompt.Suggest{
						{Text: "yes", Description: "Make this my default context"},
						{Text: "no", Description: "I will switch to this context when I need it"},
					}
					return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
				}
				fmt.Print("Make current context? ")
				answer := strings.ToLower(prompt.Input("> ", f, prompt.OptionShowCompletionAtStart()))
				if answer == "yes" || answer == "no" {
					haveMakeCurrentContextAnswer = true
					if answer == "yes" {
						context.CurrentContext = true
					}
				} else if answer == "exit" {
					userInputExitEarly <- true
					break WHILEMOREADDS
				}
			}

			// more contexts to add ?
			for !haveMoreContextsToAddAnswer {
				f := func(d prompt.Document) []prompt.Suggest {
					s := []prompt.Suggest{
						{Text: "yes", Description: "I have more auth contexts to add"},
						{Text: "no", Description: "I'm all done adding auth contexts"},
					}
					return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
				}
				fmt.Print("Add another context? ")
				answer := strings.ToLower(prompt.Input("> ", f, prompt.OptionShowCompletionAtStart()))
				if answer == "no" {
					haveMoreContextsToAddAnswer = true
					moreAdds = false
				} else if answer == "yes" {
					haveMoreContextsToAddAnswer = true
					moreAdds = true
				} else if answer == "exit" {
					userInputExitEarly <- true
					break WHILEMOREADDS
				}
			}

			// when these defaults are unsuitable, see `auth update-context` to change it later
			context.Url = "https://api.liquidweb.com"
			context.Timeout = 90
			context.Insecure = false

			validateFields := map[interface{}]interface{}{
				context.Url:         "HttpsLiquidwebUrl",
				context.Timeout:     "PositiveInt",
				context.ContextName: "NonEmptyString",
				context.Username:    "NonEmptyString",
				context.Password:    "NonEmptyString",
			}
			if err := validate.Validate(validateFields); err != nil {
				userInputError <- err
				break WHILEMOREADDS
			}

			// send context over
			userInputContext <- context
		}

		// all done
		userInputComplete <- true
	}()

	var contexts []cmdTypes.AuthContext
WAIT:
	for {
		select {
		case exitEarly := <-userInputExitEarly:
			if exitEarly {
				break WAIT
			}
		case userInputErr := <-userInputError:
			err = userInputErr
			break WAIT
		case context := <-userInputContext:
			contexts = append(contexts, context)
		case complete := <-userInputComplete:
			if complete {
				// wipe the config for a clean slate.
				if err := writeEmptyConfig(); err != nil {
					lwCliInst.Die(err)
				}
				cfgFile, cfgPathErr := getExpectedConfigPath()
				if cfgPathErr != nil {
					err = cfgPathErr
					return
				}
				if utils.FileExists(cfgFile) {
					if err := os.Remove(cfgFile); err != nil {
						lwCliInst.Die(err)
					}
					f, err := os.Create(filepath.Clean(cfgFile))
					if err != nil {
						lwCliInst.Die(err)
					}
					if err := f.Close(); err != nil {
						lwCliInst.Die(err)
					}
					if err := os.Chmod(cfgFile, 0600); err != nil {
						lwCliInst.Die(err)
					}

					if err := lwCliInst.Viper.ReadConfig(bytes.NewBuffer([]byte{})); err != nil {
						lwCliInst.Die(err)
					}
				}

				// set Viper config from contexts slice
				for _, context := range contexts {
					// ContextName
					lwCliInst.Viper.Set(fmt.Sprintf(
						"liquidweb.api.contexts.%s.contextname", context.ContextName), context.ContextName)
					// Username
					lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s.username",
						context.ContextName), context.Username)
					// Password
					lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s.password",
						context.ContextName), context.Password)
					// CurrentContext
					if context.CurrentContext {
						lwCliInst.Viper.Set("liquidweb.api.current_context", context.ContextName)
					}
					// Url
					lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s.url",
						context.ContextName), context.Url)
					// Timeout
					lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s.timeout",
						context.ContextName), context.Timeout)
					// Insecure
					lwCliInst.Viper.Set(fmt.Sprintf("liquidweb.api.contexts.%s.insecure",
						context.ContextName), context.Insecure)
				}

				// no errors or early exits, so signify to write the config just set then break
				writeConfig = true
				break WAIT
			}
		}
	}

	return
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

	f, err := os.Create(filepath.Clean(cfgFile))
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	if err := os.Chmod(cfgFile, 0600); err != nil {
		return err
	}

	return nil
}
