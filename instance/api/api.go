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
package api

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/viper"

	lwApi "github.com/liquidweb/go-lwApi"

	"github.com/liquidweb/liquidweb-cli/types/errors"
)

func New(viper *viper.Viper) (*LwCliApiClient, error) {
	if err := viper.ReadInConfig(); err != nil {
		strErr := strings.ToUpper(cast.ToString(err))
		// only consider this an error here if there is a config file
		if !strings.Contains(strErr, "CONFIG FILE") && !strings.Contains(strErr, "NOT FOUND IN") {
			return &LwCliApiClient{}, fmt.Errorf("%w Raw error: %s", errorTypes.ErrorReadingConfig, strErr)
		}
	}

	// create the object from the current context if there is one. If "auth init" has not yet been ran,
	// there would be no current context yet.
	lwCliApiClient := LwCliApiClient{Viper: viper}
	currentContext := viper.GetString("liquidweb.api.current_context")
	if currentContext != "" {
		apiUsername := viper.GetString(fmt.Sprintf("liquidweb.api.contexts.%s.username", currentContext))
		apiPassword := viper.GetString(fmt.Sprintf("liquidweb.api.contexts.%s.password", currentContext))

		lwApiCfg := lwApi.LWAPIConfig{
			Username: &apiUsername,
			Password: &apiPassword,
			Url:      viper.GetString(fmt.Sprintf("liquidweb.api.contexts.%s.url", currentContext)),
			Insecure: viper.GetBool(fmt.Sprintf("liquidweb.api.contexts.%s.insecure", currentContext)),
			Timeout: cast.ToUint(viper.GetInt(fmt.Sprintf("liquidweb.api.contexts.%s.timeout",
				currentContext))),
		}

		lwApiClient, err := lwApi.New(&lwApiCfg)
		if err != nil {
			return &LwCliApiClient{}, err
		}

		lwCliApiClient.LwApiClient = lwApiClient
	}

	return &lwCliApiClient, nil
}
