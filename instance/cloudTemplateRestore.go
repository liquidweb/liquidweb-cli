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
	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/validate"
)

type CloudTemplateRestoreParams struct {
	Template string `yaml:"template"`
	UniqId   string `yaml:"uniq-id"`
}

func (ci *Client) CloudTemplateRestore(params *CloudTemplateRestoreParams) (string, error) {
	validateFields := map[interface{}]interface{}{
		params.UniqId: "UniqId",
	}
	if err := validate.Validate(validateFields); err != nil {
		return "", err
	}

	apiArgs := map[string]interface{}{"template": params.Template, "uniq_id": params.UniqId}

	var details apiTypes.CloudTemplateRestoreResponse
	err := ci.CallLwApiInto("bleed/storm/template/restore", apiArgs, &details)
	if err != nil {
		return "", err
	}

	return details.Reimaged, nil
}
