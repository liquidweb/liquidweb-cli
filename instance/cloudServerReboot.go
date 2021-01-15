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

	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/validate"
)

type CloudServerRebootParams struct {
	UniqId string `yaml:"uniq-id"`
	Force  bool   `yaml:"force"`
}

func (self *CloudServerRebootParams) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// define defaults
	type rawType CloudServerRebootParams
	raw := rawType{} // Put your defaults here
	if err := unmarshal(&raw); err != nil {
		return err
	}
	*self = CloudServerRebootParams(raw)

	return nil
}

func (self *Client) CloudServerReboot(params *CloudServerRebootParams) (result string, err error) {

	validateFields := map[interface{}]interface{}{
		params.UniqId: "UniqId",
	}

	if err = validate.Validate(validateFields); err != nil {
		return
	}

	force := 0
	if params.Force {
		force = 1
	}

	rebootArgs := map[string]interface{}{
		"uniq_id": params.UniqId,
		"force":   force,
	}

	var resp apiTypes.CloudServerRebootResponse
	if err = self.CallLwApiInto("bleed/storm/server/reboot", rebootArgs, &resp); err != nil {
		return
	}

	result = fmt.Sprintf("Rebooted: %s\n", resp.Rebooted)

	return
}
