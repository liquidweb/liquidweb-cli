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
package validate

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cast"

	"github.com/liquidweb/liquidweb-cli/utils"
)

var ValidationFailure = errors.New("validation failed")

type InputTypes struct {
	UniqId                  InputTypeUniqId
	IP                      InputTypeIP
	PositiveInt64           InputTypePositiveInt64
	PositiveInt             InputTypePositiveInt
	NonEmptyString          InputTypeNonEmptyString
	LoadBalancerStrategy    InputTypeLoadBalancerStrategyString
	HttpsLiquidwebUrl       InputTypeHttpsLiquidwebUrl
	LoadBalancerServicePair InputTypeLoadBalancerServicePair
}

// UniqId

type InputTypeUniqId struct {
	UniqId string
}

func (x InputTypeUniqId) Validate() error {
	// must be uppercase
	allUpper := strings.ToUpper(x.UniqId)
	if allUpper != x.UniqId {
		return fmt.Errorf("a uniq_id must be uppercase")
	}

	// must be 6 characters
	if len(x.UniqId) != 6 {
		return fmt.Errorf("a uniq_id must be 6 characters long")
	}

	// must be alphanumeric
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return err
	}
	regexStr := reg.ReplaceAllString(x.UniqId, "")
	if regexStr != x.UniqId {
		return fmt.Errorf("a uniq_id must be alphanumeric")
	}

	return nil
}

// IP

type InputTypeIP struct {
	IP string
}

func (x InputTypeIP) Validate() error {

	if !utils.IpIsValid(x.IP) {
		return fmt.Errorf("ip [%s] is not a valid IP address", x.IP)
	}

	return nil
}

// PositiveInt64

type InputTypePositiveInt64 struct {
	PositiveInt64 int64
}

func (x InputTypePositiveInt64) Validate() error {
	if x.PositiveInt64 < 0 {
		return fmt.Errorf("PositiveInt64 is not > 0")
	}

	return nil
}

// PositiveInt

type InputTypePositiveInt struct {
	PositiveInt int
}

func (x InputTypePositiveInt) Validate() error {
	if x.PositiveInt < 0 {
		return fmt.Errorf("PositiveInt is not > 0")
	}

	return nil
}

// NonEmptyString

type InputTypeNonEmptyString struct {
	NonEmptyString string
}

func (x InputTypeNonEmptyString) Validate() error {
	if x.NonEmptyString == "" {
		return fmt.Errorf("NonEmptyString cannot be empty")
	}

	return nil
}

// LoadBalancerStrategy

type InputTypeLoadBalancerStrategyString struct {
	LoadBalancerStrategy string
}

func (x InputTypeLoadBalancerStrategyString) Validate() error {
	strategies := map[string]int{
		"roundrobin":  1,
		"connections": 1,
		"cells":       1,
	}

	if _, exists := strategies[x.LoadBalancerStrategy]; !exists {
		var slice []string
		slice = append(slice, fmt.Sprintf("LoadBalancer strategy [%s] is invalid. Valid strategies: ",
			x.LoadBalancerStrategy))
		for strategy, _ := range strategies {
			slice = append(slice, fmt.Sprintf("%s ", strategy))
		}
		return fmt.Errorf("%s", strings.Join(slice[:], ""))
	}

	return nil
}

// HttpsLiquidwebUrl

type InputTypeHttpsLiquidwebUrl struct {
	HttpsLiquidwebUrl string
}

func (x InputTypeHttpsLiquidwebUrl) Validate() error {
	if !strings.HasPrefix(x.HttpsLiquidwebUrl, "https://") {
		return fmt.Errorf("given url [%s] appears invalid; should start with 'https://'",
			x.HttpsLiquidwebUrl)
	}

	if !strings.Contains(x.HttpsLiquidwebUrl, "liquidweb.com") {
		return fmt.Errorf("given url [%s] appears invalid; should contain 'liquidweb.com'",
			x.HttpsLiquidwebUrl)
	}

	if _, err := url.ParseRequestURI(x.HttpsLiquidwebUrl); err != nil {
		return fmt.Errorf("given url [%s] appears invalid; %s", x.HttpsLiquidwebUrl, err)
	}

	return nil
}

// LoadBalancerServicePair

type InputTypeLoadBalancerServicePair struct {
	LoadBalancerServicePair string
}

func (x InputTypeLoadBalancerServicePair) Validate() error {
	if !strings.Contains(x.LoadBalancerServicePair, ":") {
		return fmt.Errorf("given LoadBalancerServicePair [%s] contains no ':' which is invalid",
			x.LoadBalancerServicePair)
	}

	splitPair := strings.Split(x.LoadBalancerServicePair, ":")

	if len(splitPair) != 2 {
		return fmt.Errorf(
			"A LoadBalancerServicePair must contain exactly one source/destination port pair")
	}

	for _, portStr := range splitPair {
		if _, err := strconv.Atoi(portStr); err != nil {
			return fmt.Errorf("port [%s] in port pair [%s] doesnt look numeric", portStr,
				x.LoadBalancerServicePair)
		}

		portInt := cast.ToInt(portStr)
		if portInt <= 0 || portInt > 65535 {
			return fmt.Errorf("port [%d] is invalid; must be between 1 and 65535", portInt)
		}
	}

	return nil
}
