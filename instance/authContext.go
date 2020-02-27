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
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"

	"github.com/liquidweb/liquidweb-cli/types/cmd"
)

func (client *Client) RemoveContext(context string) error {
	// review this function once Unset (or similar) in viper is available
	// https://github.com/spf13/viper/pull/519
	if context == "" {
		return fmt.Errorf("context cannot be empty")
	}

	currentContext := client.Viper.GetString("liquidweb.api.current_context")
	if context == currentContext {
		return fmt.Errorf("cannot remove context currently set as current context")
	}

	contexts := client.Viper.GetStringMap("liquidweb.api.contexts")
	if _, exists := contexts[context]; !exists {
		return fmt.Errorf("context %s doesnt exist, cannot remove", context)
	}

	// save current config into a map, then delete requested context
	cfgMap := client.Viper.AllSettings()
	delete(cfgMap["liquidweb"].(map[string]interface{})["api"].(map[string]interface{})["contexts"].(map[string]interface{}), context)

	// json encode modified map
	encodedCfg, err := json.MarshalIndent(cfgMap, "", " ")
	if err != nil {
		return err
	}

	// read newly encoded config back into viper
	if err := client.Viper.ReadConfig(bytes.NewBuffer(encodedCfg)); err != nil {
		return err
	}

	// write the new viper configuration to file
	if err := client.Viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func ValidateContext(wantedContext string, vp *viper.Viper) error {
	var isValid bool
	contexts := vp.GetStringMap("liquidweb.api.contexts")
	for _, contextInter := range contexts {
		var context cmdTypes.AuthContext
		if err := CastFieldTypes(contextInter, &context); err != nil {
			return err
		}

		if context.ContextName == wantedContext {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("given context [%s] is not a valid context", wantedContext)
	}

	return nil
}
