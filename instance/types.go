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
	"github.com/spf13/viper"

	lwApi "github.com/liquidweb/go-lwApi"
)

type Client struct {
	LwApiClient *lwApi.Client
	Viper       *viper.Viper
}

type AllPaginatedResultsArgs struct {
	Method         string
	MethodArgs     map[string]interface{}
	ResultsPerPage int64
}