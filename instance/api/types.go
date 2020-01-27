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
	"github.com/spf13/viper"

	lwApi "github.com/liquidweb/go-lwApi"

	"github.com/liquidweb/liquidweb-cli/types/errors"
)

type LwCliApiClient struct {
	LwApiClient *lwApi.Client
	Viper       *viper.Viper
}

func (x LwCliApiClient) Call(method string, params interface{}) (got interface{}, err error) {
	if err = x.Viper.ReadInConfig(); err != nil {
		err = fmt.Errorf("%w Raw error: %s", errorTypes.InvalidConfigSyntax, err)
		return
	}

	currentContext := x.Viper.GetString("liquidweb.api.current_context")
	if currentContext == "" {
		err = errorTypes.NoCurrentContext
		return
	}

	got, err = x.LwApiClient.Call(method, params)

	return
}
