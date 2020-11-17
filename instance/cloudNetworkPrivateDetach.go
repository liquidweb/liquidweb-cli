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

type CloudNetworkPrivateDetachParams struct {
	UniqId []string `yaml:"uniq-id"`
}

func (self *CloudNetworkPrivateDetachParams) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// define defaults
	type rawType CloudNetworkPrivateDetachParams
	raw := rawType{} // Put your defaults here
	if err := unmarshal(&raw); err != nil {
		return err
	}
	*self = CloudNetworkPrivateDetachParams(raw)

	return nil
}

func (self *Client) CloudNetworkPrivateDetach(params *CloudNetworkPrivateDetachParams) (result string, err error) {
	if len(params.UniqId) == 0 {
		err = errors.New("--uniq-id must be given")
		return
	}

	var b bytes.Buffer

	for _, uniqId := range params.UniqId {
		validateFields := map[interface{}]interface{}{
			uniqId: "UniqId",
		}
		if err := validate.Validate(validateFields); err != nil {
			fmt.Printf("uniqId [%s] is invalid; ignoring...\n", uniqId)
			continue
		}

		apiArgs := map[string]interface{}{"uniq_id": uniqId}

		var attachedDetails apiTypes.CloudNetworkPrivateIsAttachedResponse
		if err = self.CallLwApiInto("bleed/network/private/isattached", apiArgs, &attachedDetails); err != nil {
			return
		}
		if !attachedDetails.IsAttached {
			err = errors.New("Cloud Server is already detached from the Private Network")
			return
		}

		var details apiTypes.CloudNetworkPrivateDetachResponse
		if err = self.CallLwApiInto("bleed/network/private/detach", apiArgs, &details); err != nil {
			return
		}

		b.WriteString(fmt.Sprintf("Detaching %s from private network\n", details.Detached))
		b.WriteString(fmt.Sprintf("\n\nYou can check progress with 'cloud server status --uniq-id %s'\n", uniqId))
	}

	result = b.String()

	return
}
