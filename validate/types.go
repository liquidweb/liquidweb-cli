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
	"fmt"
	"strings"

	"github.com/liquidweb/liquidweb-cli/utils"
)

type InputTypes struct {
	UniqId         InputTypeUniqId
	IP             InputTypeIP
	PositiveInt64  InputTypePositiveInt64
	PostiveInt     InputTypePositiveInt
	NonEmptyString InputTypeNonEmptyString
}

// UniqId

type InputTypeUniqId struct {
	UniqId string
}

func (x InputTypeUniqId) Validate() error {
	allUpper := strings.ToUpper(x.UniqId)
	if allUpper != x.UniqId {
		return fmt.Errorf("a uniq_id must be uppercase")
	}

	if len(x.UniqId) != 6 {
		return fmt.Errorf("a uniq_id must be 6 characters long")
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
