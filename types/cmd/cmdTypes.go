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
package cmdTypes

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cast"

	"github.com/liquidweb/liquidweb-cli/validate"
)

type AuthContext struct {
	CurrentContext bool   `json:"currentcontext" mapstructure:"currentcontext"`
	ContextName    string `json:"contextname" mapstructure:"contextname"`
	Username       string `json:"username" mapstructure:"username"`
	Password       string `json:"password" mapstructure:"password"`
	Url            string `json:"url" mapstructure:"url"`
	Insecure       bool   `json:"insecure" mapstructure:"insecure"`
	Timeout        int    `json:"timeout" mapstructure:"timeout"`
}

type LoadBalancerHealthCheck struct {
	HealthCheck map[string]string `json:"health_check" mapstructure:"health_check"`
}

func (x LoadBalancerHealthCheck) Transform() (bySrcPort map[string]map[string]interface{}, err error) {
	// this is a bit of a hack. I couldn't find any supported way in cobra/viper/pflag to
	// get flags to go into a slice of maps. So its reading from flags into a single map,
	// and then doing this logic here to transform that into one suitable to aid later
	// consumption by the api.

	bySrcPort = make(map[string]map[string]interface{})

	for key, value := range x.HealthCheck {
		re := regexp.MustCompile(`^\d+`)
		srcPort := cast.ToString(re.Find([]byte(key)))
		re = regexp.MustCompile(`^\d+_(.*)`)
		strippedKeyMatch := re.FindStringSubmatch(key)
		if len(strippedKeyMatch) != 2 {
			err = fmt.Errorf("error parsing health-check flags; flag [%s] isn't of expected format [%+v]. Make sure you placed $port_ infront of parameter name.", key, value)
			return
		}
		param := strippedKeyMatch[1]

		// zero value of a map is nil. If the map isn't yet initialized, its now time to do so
		if bySrcPort[srcPort] == nil {
			bySrcPort[srcPort] = make(map[string]interface{})
		}

		// http_body_match
		if param == "http_body_match" {
			bySrcPort[srcPort]["http_body_match"] = value
		}

		// http_path
		if param == "http_path" {
			bySrcPort[srcPort]["http_path"] = value
		}

		// http_use_tls
		if param == "http_use_tls" {
			var boolValue bool
			if value == "true" {
				boolValue = true
			}
			bySrcPort[srcPort]["http_use_tls"] = boolValue
		}

		// http_response_codes
		if param == "http_response_codes" {
			// this is supposed to be a string delimited by a ','... '200-204,300-301,404'
			bySrcPort[srcPort]["http_response_codes"] = strings.ReplaceAll(value, ":", ",")
		}

		// timeout
		// when not set, api defaults to 5
		if param == "timeout" {
			bySrcPort[srcPort]["timeout"] = value
		}

		// failure_threshold
		// when not set, api defaults to 3
		if param == "failure_threshold" {
			bySrcPort[srcPort]["failure_threshold"] = value
		}

		// protocol
		if param == "protocol" {
			bySrcPort[srcPort]["protocol"] = value
		}

		// interval
		// when not set, api defaults to 30
		if param == "interval" {
			bySrcPort[srcPort]["interval"] = value
		}
	}

	// bySrcPort is now built, verify the inputs
	for sourcePort, healthCheck := range bySrcPort {
		// protocol is required
		if _, exists := healthCheck["protocol"]; !exists {
			err = fmt.Errorf("protocol is required and was not given for service with source port [%+v]", sourcePort)
			return
		}

		// place defaults for http_path, http_use_tls, http_response_codes if protocol == "http" if unset.
		if healthCheck["protocol"] == "http" {
			// if http_path wasn't passed, default to /
			if _, exists := healthCheck["http_path"]; !exists {
				bySrcPort[sourcePort]["http_path"] = "/"
			}
			// if http_response_codes wasn't passed, default to '200-206,300-304'
			if _, exists := healthCheck["http_response_codes"]; !exists {
				bySrcPort[sourcePort]["http_response_codes"] = "200-206,300-304"
			}
			// if http_use_tls wasn't passed, default it to false
			if _, exists := healthCheck["http_use_tls"]; !exists {
				bySrcPort[sourcePort]["http_use_tls"] = false
			}
		} else {
			// when protocol isn't http, these shouldn't be set.
			if _, exists := healthCheck["http_path"]; exists {
				err = fmt.Errorf("http_path cannot be set when protocol isn't http")
				return
			}
			if _, exists := healthCheck["http_response_codes"]; exists {
				err = fmt.Errorf("http_response_codes cannot be set when protocol isn't http")
				return
			}
			if _, exists := healthCheck["http_use_tls"]; exists {
				err = fmt.Errorf("http_use_tls cannot be set when protocol isn't http")
				return
			}
			if _, exists := healthCheck["http_body_match"]; exists {
				err = fmt.Errorf("http_body_match cannot be set when protocol isn't http")
				return
			}
		}

		validateFields := map[interface{}]interface{}{
			healthCheck["protocol"]: "LoadBalancerHealthCheckProtocol",
		}

		if val, exists := healthCheck["http_response_codes"]; exists {
			validateFields[val] = "LoadBalancerHttpCodeRange"
		}
		if val, exists := healthCheck["timeout"]; exists {
			if _, convErr := strconv.Atoi(cast.ToString(val)); convErr != nil {
				err = fmt.Errorf("timeout value [%+v] doesn't look numeric", val)
				return
			}
			validateFields[cast.ToInt(val)] = "PositiveInt"
		}
		if val, exists := healthCheck["interval"]; exists {
			if _, convErr := strconv.Atoi(cast.ToString(val)); convErr != nil {
				err = fmt.Errorf("interval value [%+v] doesn't look numeric", val)
				return
			}
			validateFields[cast.ToInt(val)] = "PositiveInt"
		}
		if val, exists := healthCheck["failure_threshold"]; exists {
			if _, convErr := strconv.Atoi(cast.ToString(val)); convErr != nil {
				err = fmt.Errorf("failure_threshold value [%+v] doesn't look numeric", val)
				return
			}
			validateFields[cast.ToInt(val)] = "PositiveInt"
		}

		if validateErr := validate.Validate(validateFields); validateErr != nil {
			err = fmt.Errorf("healthCheck validation failed for service with source port [%+v]: %s", sourcePort, validateErr)
			return
		}
	}

	return
}
