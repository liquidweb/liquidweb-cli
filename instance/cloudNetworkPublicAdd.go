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
	"fmt"

	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/validate"
)

type CloudNetworkPublicAddParams struct {
	UniqId       string   `yaml:"uniq-id"`
	ConfigureIps bool     `yaml:"configure-ips"`
	NewIps       int64    `yaml:"new-ips"`
	NewIp6s      int64    `yaml:"new-ip6s"`
	PoolIps      []string `yaml:"pool-ips"`
	Pool6Ips     []string `yaml:"pool6-ips"`
}

func (self *CloudNetworkPublicAddParams) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// define defaults
	type rawType CloudNetworkPublicAddParams
	raw := rawType{} // Put your defaults here
	if err := unmarshal(&raw); err != nil {
		return err
	}
	*self = CloudNetworkPublicAddParams(raw)

	return nil
}

func (self *Client) CloudNetworkPublicAdd(params *CloudNetworkPublicAddParams) (result string, err error) {
	validateFields := map[interface{}]interface{}{
		params.UniqId: "UniqId",
	}
	if err = validate.Validate(validateFields); err != nil {
		return
	}

	if params.NewIps == 0 && len(params.PoolIps) == 0 && params.NewIp6s == 0 && len(params.Pool6Ips) == 0 {
		err = fmt.Errorf("at least one of --new-ips --pool-ips --new-ip6s --pool6-ips must be given")
		return
	}

	apiArgs := map[string]interface{}{
		"configure_ips": params.ConfigureIps,
		"uniq_id":       params.UniqId,
	}
	if params.NewIps != 0 {
		apiArgs["ip_count"] = params.NewIps
		validateFields := map[interface{}]interface{}{params.NewIps: "PositiveInt64"}
		if err = validate.Validate(validateFields); err != nil {
			return
		}
	}
	if len(params.PoolIps) != 0 {
		apiArgs["pool_ips"] = params.PoolIps
		validateFields := map[interface{}]interface{}{}
		for _, ip := range params.PoolIps {
			validateFields[ip] = "IP"
		}
		if err = validate.Validate(validateFields); err != nil {
			return
		}
	}

	if params.NewIp6s != 0 {
		apiArgs["ip6_count"] = params.NewIp6s
		validateFields := map[interface{}]interface{}{params.NewIps: "PositiveInt64"}
		if err = validate.Validate(validateFields); err != nil {
			return
		}
	}
	if len(params.Pool6Ips) != 0 {
		apiArgs["pool6_ips"] = params.Pool6Ips
		validateFields := map[interface{}]interface{}{}
		for _, ip := range params.Pool6Ips {
			validateFields[ip] = "CIDR"
		}
		if err = validate.Validate(validateFields); err != nil {
			return
		}
	}

	var details apiTypes.NetworkIpAdd
	if err = self.CallLwApiInto("bleed/network/ip/add", apiArgs, &details); err != nil {
		return
	}

	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("Adding [%s] to Cloud Server\n", details.Adding))

	if params.ConfigureIps {
		b.WriteString(fmt.Sprint("IP(s) will be automatically configured.\n"))
	} else {
		b.WriteString(fmt.Sprint("IP(s) will need to be manually configured.\n"))
	}

	result = b.String()

	return
}
