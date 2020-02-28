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
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/liquidweb/liquidweb-cli/instance"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Process YAML plan file",
	Long: `Process YAML plan file.

Examples:
'lw plan --file plan.yaml'

Example plan file to create a cloud server:

---
cloud:
   server:
      create:
         - type: "SS.VPS"
           template: "UBUNTU_1804_UNMANAGED"
           zone: 40460
           hostname: "db1.somedomain.com"
           ips: 1
           public-ssh-key: "public ssh key string here "
           config_id: 88
         - type: "SS.VPS"
           template: "UBUNTU_1804_UNMANAGED"
           zone: 40460
           hostname: "web1.somedomain.com"
           ips: 1
           public-ssh-key: "public ssh key string here "
           config_id: 88

`,
	Run: func(cmd *cobra.Command, args []string) {
		planFile, _ := cmd.Flags().GetString("file")
		varFlag, err := cmd.Flags().GetStringSlice("var")

		if err != nil {
			panic(err)
		}

		_, err = os.Stat(planFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Plan file \"%s\" does not exist.\n", planFile)
				os.Exit(1)
			} else {
				panic(err)
			}
		}

		planYaml, err := ioutil.ReadFile(planFile)
		if err != nil {
			panic(err)
		}

		if len(varFlag) > 0 {
			vars := make(map[string]string)
			for _, v := range varFlag {
				s := strings.Split(v, "=")
				vars[s[0]] = s[1]
			}

			var tmplBytes bytes.Buffer
			tmpl, err := template.New("plan.yaml").Parse(string(planYaml))
			if err != nil {
				panic(err)
			}
			err = tmpl.Execute(&tmplBytes, vars)
			if err != nil {
				panic(err)
			}
			planYaml = tmplBytes.Bytes()
		}

		var plan instance.Plan
		err = yaml.Unmarshal(planYaml, &plan)
		if err != nil {
			fmt.Printf("Error parsing YAML file: %s\n", err)
		}

		if err := lwCliInst.ProcessPlan(&plan); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)

	planCmd.Flags().String("file", "", "YAML file used to define a plan")
	planCmd.Flags().StringSlice("var", nil, "define variable name")
	planCmd.MarkFlagRequired("file")
}
