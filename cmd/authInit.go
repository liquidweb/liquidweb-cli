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
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/types/cmd"
	"github.com/liquidweb/liquidweb-cli/types/errors"
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
	if !terminal.IsTerminal(0) || !terminal.IsTerminal(1) {
		err = errorTypes.UnknownTerminal
		return
	}
	oldState, termMakeErr := terminal.MakeRaw(0)
	if termMakeErr != nil {
		err = termMakeErr
		return
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
		if _, err = term.Write([]byte("Warning: This will delete all auth contexts. Continue (yes/[no])?: ")); err != nil {
			return
		}
		proceedBytes, readErr := term.ReadLine()
		if readErr != nil {
			err = readErr
			return
		}
		proceedString := cast.ToString(proceedBytes)
		if proceedString != "yes" && proceedString != "no" && proceedString != "" {
			if _, err = term.Write([]byte("invalid input.\n")); err != nil {
				return
			}
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
		return
	}

	// create context add loop channels
	userInputComplete := make(chan bool)
	userInputError := make(chan error)
	userInputExitEarly := make(chan bool)
	userInputContext := make(chan cmdTypes.AuthContext)

	if _, err = term.Write([]byte("\nTo exit early, type 'exit' or send EOF (ctrl+d)\n\n")); err != nil {
		return
	}

	// start context add loop
	go func() {
	WHILEMOREADDS:
		for moreAdds {
			var (
				context                      cmdTypes.AuthContext
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
				if _, err := term.Write([]byte("Name this context: ")); err != nil {
					userInputError <- err
					break WHILEMOREADDS
				}
				contextNameBytes, err := term.ReadLine()
				if err != nil {
					userInputError <- err
					break WHILEMOREADDS
				}
				contextNameAnswer = cast.ToString(contextNameBytes)
				if contextNameAnswer == "exit" {
					userInputExitEarly <- true
					break WHILEMOREADDS
				} else if contextNameAnswer == "" {
					if _, err := term.Write([]byte("context name cannot be blank.\n")); err != nil {
						userInputError <- err
						break WHILEMOREADDS
					}
				} else if strings.ToLower(contextNameAnswer) != contextNameAnswer {
					if _, err := term.Write([]byte("context name cannot be uppercase.\n")); err != nil {
						userInputError <- err
						break WHILEMOREADDS
					}
				} else {
					haveContextNameAnswer = true
					context.ContextName = contextNameAnswer
				}
			}

			// username
			for !haveUsernameAnswer {
				if _, err := term.Write([]byte("LiquidWeb username: ")); err != nil {
					userInputError <- err
					break WHILEMOREADDS
				}
				usernameBytes, err := term.ReadLine()
				if err != nil {
					userInputError <- err
					break WHILEMOREADDS
				}
				usernameAnswer = cast.ToString(usernameBytes)
				if usernameAnswer == "exit" {
					userInputExitEarly <- true
					break WHILEMOREADDS
				} else if usernameAnswer == "" {
					if _, err := term.Write([]byte("username cannot be blank.\n")); err != nil {
						userInputError <- err
						break WHILEMOREADDS
					}
				} else {
					haveUsernameAnswer = true
					context.Username = usernameAnswer
				}
			}

			// password
			for !havePasswordAnswer {
				passwordBytes, err := term.ReadPassword("LiquidWeb password: ")
				if err != nil {
					userInputError <- err
					break WHILEMOREADDS
				}
				passwordAnswer = cast.ToString(passwordBytes)
				if passwordAnswer == "exit" {
					userInputExitEarly <- true
					break WHILEMOREADDS
				} else if passwordAnswer == "" {
					if _, err := term.Write([]byte("password cannot be blank.\n")); err != nil {
						userInputError <- err
						break WHILEMOREADDS
					}
				} else {
					havePasswordAnswer = true
					context.Password = passwordAnswer
				}
			}

			// make current context?
			for !haveMakeCurrentContextAnswer {
				if _, err := term.Write([]byte("Make current context? ([yes]/no)")); err != nil {
					userInputError <- err
					break WHILEMOREADDS
				}
				makeCurrentContextBytes, err := term.ReadLine()
				if err != nil {
					userInputError <- err
					break WHILEMOREADDS
				}
				makeCurrentContextString := cast.ToString(makeCurrentContextBytes)
				if makeCurrentContextString == "exit" {
					userInputExitEarly <- true
					break WHILEMOREADDS
				}
				if makeCurrentContextString != "" && makeCurrentContextString != "yes" && makeCurrentContextString != "no" {
					if _, err := term.Write([]byte("invalid input.\n")); err != nil {
						userInputError <- err
						break WHILEMOREADDS
					}
					continue
				}

				haveMakeCurrentContextAnswer = true
				if makeCurrentContextString == "yes" || makeCurrentContextString == "" {
					context.CurrentContext = true
				}
			}

			// more contexts to add ?
			for !haveMoreContextsToAddAnswer {
				if _, err := term.Write([]byte("Add another context? (yes/[no]): ")); err != nil {
					userInputError <- err
					break WHILEMOREADDS
				}
				moreContextsBytes, err := term.ReadLine()
				if err != nil {
					userInputError <- err
					break WHILEMOREADDS
				}

				answer := cast.ToString(moreContextsBytes)
				if answer == "exit" {
					userInputExitEarly <- true
					break WHILEMOREADDS
				}
				if answer != "" && answer != "yes" && answer != "no" {
					if _, err := term.Write([]byte("invalid input.\n")); err != nil {
						userInputError <- err
						break WHILEMOREADDS
					}
					continue
				}

				if answer == "no" || answer == "" {
					moreAdds = false
				}

				haveMoreContextsToAddAnswer = true
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
					f, err := os.Create(cfgFile)
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

	f, err := os.Create(cfgFile)
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
