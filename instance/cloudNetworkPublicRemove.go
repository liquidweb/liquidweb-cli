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
	"bytes"
	"errors"
	"fmt"

	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/validate"
)

type CloudNetworkPublicRemoveParams struct {
	UniqId       string   `yaml:"uniq-id"`
	ConfigureIps bool     `yaml:"configure-ips"`
	Ips          []string `yaml:"ips"`
}

func (self *CloudNetworkPublicRemoveParams) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// define defaults
	type rawType CloudNetworkPublicRemoveParams
	raw := rawType{} // Put your defaults here
	if err := unmarshal(&raw); err != nil {
		return err
	}
	*self = CloudNetworkPublicRemoveParams(raw)

	return nil
}

func (self *Client) CloudNetworkPublicRemove(params *CloudNetworkPublicRemoveParams) (result string, err error) {
	if len(params.UniqId) == 0 {
		err = errors.New("--uniq-id is required")
	}

	validateFields := map[interface{}]interface{}{
		params.UniqId: "UniqId",
	}
	if err = validate.Validate(validateFields); err != nil {
		return
	}

	apiArgs := map[string]interface{}{
		"configure_ips": params.ConfigureIps,
		"uniq_id":       params.UniqId,
	}

	var b bytes.Buffer
	for _, ip := range params.Ips {
		validateFields := map[interface{}]interface{}{
			ip: "IpOrCidr",
		}
		if err := validate.Validate(validateFields); err != nil {
			fmt.Printf("%s ... skipping\n", err)
			continue
		}

		var details apiTypes.NetworkIpRemove
		apiArgs["ip"] = ip
		if err = self.CallLwApiInto("bleed/network/ip/remove", apiArgs, &details); err != nil {
			return
		}

		b.WriteString(fmt.Sprintf("Removing [%s] from Cloud Server\n", details.Removing))

		if params.ConfigureIps {
			b.WriteString(fmt.Sprint("IP(s) will be automatically removed from the network configuration.\n"))
		} else {
			b.WriteString(fmt.Sprint("IP(s) will need to be manually removed from the network configuration.\n"))
		}
	}

	result = b.String()

	return
}
