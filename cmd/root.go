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
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/types/api"
)

var cfgFile string
var lwCliInst instance.Client

var rootCmd = &cobra.Command{
	Use:   "liquidweb-cli",
	Short: "CLI interface for LiquidWeb",
	Long: `CLI interface for LiquidWeb.

Command line interface for interacting with LiquidWeb services via
LiquidWebs Public API.

If this is your first time running, you will need to setup auth
contexts. An auth context contains authenication data for accessing
your LiquidWeb Account. To setup your first auth context, you can
run 'auth init'. For further information on auth contexts, be sure
to checkout 'help auth' for a list of capabilities.

Consult the various subcommands for specific features and
capabilities.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.liquidweb-cli.yaml)")
}

func initConfig() {
	vp := viper.New()
	if cfgFile != "" {
		// Use config file from the flag.
		vp.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".liquidweb-cli" (without extension).
		vp.AddConfigPath(home)
		vp.SetConfigName(".liquidweb-cli")
	}

	vp.AutomaticEnv()
	vp.ReadInConfig()

	var lwCliInstErr error
	lwCliInst, lwCliInstErr = instance.New(vp)
	if lwCliInstErr != nil {
		fmt.Printf("Fatal: [%s]\n", lwCliInstErr)
		os.Exit(1)
	}
}

func derivePrivateParentUniqId(name string) (string, error) {
	var (
		privateParentUniqId     string
		privateParentDetails    apiTypes.CloudPrivateParentDetails
		privateParentDetailsErr error
	)

	// if name looks like a uniq_id, try it as a uniq_id first.
	if len(name) == 6 && strings.ToUpper(name) == name {
		if err := lwCliInst.CallLwApiInto("bleed/storm/private/parent/details",
			map[string]interface{}{"uniq_id": name},
			&privateParentDetails); err == nil {
			privateParentUniqId = name
		} else {
			privateParentDetailsErr = fmt.Errorf(
				"failed fetching parent details treating given --private-parent arg as a uniq_id [%s]: %s",
				name, err)
		}
	}

	// if we havent found the pp details yet, try assuming name is the name of the pp
	if privateParentUniqId == "" {
		methodArgs := instance.AllPaginatedResultsArgs{
			Method:         "bleed/storm/private/parent/list",
			ResultsPerPage: 100,
		}
		results, err := lwCliInst.AllPaginatedResults(&methodArgs)
		if err != nil {
			lwCliInst.Die(err)
		}

		for _, item := range results.Items {
			var privateParentDetails apiTypes.CloudPrivateParentDetails
			if err := instance.CastFieldTypes(item, &privateParentDetails); err != nil {
				lwCliInst.Die(err)
			}

			if privateParentDetails.Domain == name {
				// found it get details
				err := lwCliInst.CallLwApiInto("bleed/storm/private/parent/details",
					map[string]interface{}{
						"uniq_id": privateParentDetails.UniqId,
					},
					&privateParentDetails)
				if err != nil {
					privateParentDetailsErr = fmt.Errorf(
						"failed fetching private parent details for discovered uniq_id [%s] error: %s %w",
						privateParentDetails.UniqId, err, privateParentDetailsErr)
					return "", privateParentDetailsErr
				}
				privateParentUniqId = privateParentDetails.UniqId
				break // found the uniq_id so break
			}
		}
	}

	if privateParentUniqId == "" {
		return "", fmt.Errorf("failed deriving uniq_id of private parent from [%s]: %s", name, privateParentDetailsErr)
	}

	return privateParentUniqId, nil
}
