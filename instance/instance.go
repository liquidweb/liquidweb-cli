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
	"os"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	lwCliInstApi "github.com/liquidweb/liquidweb-cli/instance/api"

	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/types/errors"
	"github.com/liquidweb/liquidweb-cli/utils"
)

func New(viper *viper.Viper) (Client, error) {

	lwApiClient, err := lwCliInstApi.New(viper)
	if err != nil {
		return Client{}, fmt.Errorf(
			"Failed creating an lwApi client. Error was:\n%s\nPlease check your liquidweb-cli config file for errors or ommissions\n",
			err)
	}

	client := Client{
		LwApiClient: lwApiClient,
		Viper:       viper,
	}

	return client, nil
}

func (*Client) Die(err error) {
	utils.PrintRed("A fatal error has occurred:\n\n")
	fmt.Printf("%s\n\nIf you need assistance please contact support. ", err)
	os.Exit(1)
}

func (*Client) JsonEncodeAndPrettyPrint(data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "    ")

	err := encoder.Encode(data)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (*Client) JsonPrettyPrint(inJson string) (string, error) {
	var outJson bytes.Buffer
	err := json.Indent(&outJson, []byte(inJson), "", "    ")
	if err != nil {
		return "", err
	}
	return outJson.String(), nil
}

func (client *Client) RemoveContext(context string) error {
	// this function should be deleted and Unset in viper used instead once
	// https://github.com/spf13/viper/pull/519
	// is merged or similar functionality is added.
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

func (client *Client) CallLwApiInto(method string, methodArgs map[string]interface{}, obj interface{}) (err error) {
	got, err := client.LwApiClient.Call(method, methodArgs)
	if err != nil {
		return
	}

	err = CastFieldTypes(got, &obj)

	return
}

func (client *Client) AllPaginatedResults(args *AllPaginatedResultsArgs) (apiTypes.MergedPaginatedList, error) {

	if args.Method == "" {
		return apiTypes.MergedPaginatedList{}, fmt.Errorf("%w Method", errorTypes.LwCliInputError)
	}

	resultsPerPage := int64(500)
	if args.ResultsPerPage != 0 {
		resultsPerPage = args.ResultsPerPage
	}

	methodArgs := args.MethodArgs
	if methodArgs == nil {
		methodArgs = map[string]interface{}{
			"page_size": resultsPerPage,
		}
	} else {
		methodArgs["page_size"] = resultsPerPage
	}

	got, err := client.LwApiClient.Call(args.Method, methodArgs)
	if err != nil {
		return apiTypes.MergedPaginatedList{}, err
	}

	var list apiTypes.PaginatedList
	if err := CastFieldTypes(got, &list); err != nil {
		return apiTypes.MergedPaginatedList{}, err
	}

	mergedList := apiTypes.MergedPaginatedList{
		Items:    list.Items,
		PageSize: resultsPerPage,
	}

	nextPage := list.PageNum + 1
	if list.PageNum < list.PageTotal {
		morePages := true

		for morePages {
			methodArgs["page_num"] = nextPage
			got, err := client.LwApiClient.Call(args.Method, methodArgs)
			if err != nil {
				return apiTypes.MergedPaginatedList{}, err
			}

			var page apiTypes.PaginatedList
			if err := CastFieldTypes(got, &page); err != nil {
				return apiTypes.MergedPaginatedList{}, err
			}

			// append page to mergedList
			for _, item := range page.Items {
				mergedList.Items = append(mergedList.Items, item)
			}

			nextPage++
			if nextPage > page.PageTotal {
				morePages = false
			}
		}
	}

	mergedList.MergedPages = nextPage - 1

	return mergedList, nil
}

func CastFieldTypes(source interface{}, dest interface{}) (err error) {
	defer func() {
		if paniced := recover(); paniced != nil {
			err = fmt.Errorf("%w source [%+v] dest type [%s]: %+v",
				errorTypes.LwApiUnexpectedResponseStructure, source,
				reflect.TypeOf(dest).String(), paniced)
		}
	}()

	if err = mapstructure.WeakDecode(source, &dest); err != nil {
		err = fmt.Errorf("%w\nsource [%+v] dest type [%s] error: %+v",
			errorTypes.LwApiUnexpectedResponseStructure, source,
			reflect.TypeOf(dest).String(), err)
	}

	return
}
