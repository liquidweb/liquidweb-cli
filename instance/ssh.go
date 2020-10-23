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
package instance

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/validate"
)

type SshParams struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	PrivateKey      string `yaml:"private-key"`
	User            string `yaml:"user"`
	AgentForwarding bool   `yaml:"agent-forwarding"`
	Command         string `yaml:"command"`
}

func (self *SshParams) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawType SshParams
	raw := rawType{
		Port: 22,
		User: "root",
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}
	*self = SshParams(raw)

	return nil
}

func (self *SshParams) TranslateHost(ci *Client) (ip string, err error) {
	validateFields := map[interface{}]interface{}{
		self.Host: "UniqId",
	}

	if err = validate.Validate(validateFields); err == nil {
		var subaccnt apiTypes.Subaccnt
		apiArgs := map[string]interface{}{
			"uniq_id": self.Host,
		}
		if aErr := ci.CallLwApiInto("bleed/asset/details", apiArgs, &subaccnt); aErr != nil {
			err = aErr
			return
		}

		ip = subaccnt.Ip
	} else {
		methodArgs := AllPaginatedResultsArgs{
			Method:         "bleed/asset/list",
			ResultsPerPage: 100,
		}
		results, aErr := ci.AllPaginatedResults(&methodArgs)
		if aErr != nil {
			err = aErr
			return
		}
		for _, item := range results.Items {
			var subaccnt apiTypes.Subaccnt
			if err = CastFieldTypes(item, &subaccnt); err != nil {
				return
			}

			if subaccnt.Domain == self.Host {
				ip = subaccnt.Ip
				break
			}
		}
	}

	if ip == "" {
		err = fmt.Errorf("unable to determine ip for Host [%s]", self.Host)
		return
	}

	return
}

func (self *Client) Ssh(params *SshParams) (err error) {
	validateFields := map[interface{}]interface{}{
		params.Port: "PositiveInt",
	}
	if err = validate.Validate(validateFields); err != nil {
		return
	}

	ip, err := params.TranslateHost(self)
	if err != nil {
		return
	}

	sshArgs := []string{}
	if params.PrivateKey != "" {
		sshArgs = append(sshArgs, "-i", params.PrivateKey)
	}

	if params.AgentForwarding {
		sshArgs = append(sshArgs, "-A")
	}

	sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", params.User, ip))
	sshArgs = append(sshArgs, fmt.Sprintf("-p %d", params.Port))

	if params.Command != "" {
		sshArgs = append(sshArgs, params.Command)
	}

	cmd := exec.Command("ssh", sshArgs...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	err = cmd.Run()

	return
}
