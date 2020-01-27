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
package errorTypes

import (
	"errors"
)

var LwCliInputError = errors.New("Invalid input; missing required paramater")
var LwApiUnexpectedResponseStructure = errors.New("Unexpected API response structure when calling method")
var UnknownTerminal = errors.New("unknown terminal")
var MergeConfigError = errors.New("error merging configuration")
var InvalidConfigSyntax = errors.New("configuration contains invalid syntax; use 'auth init' to create a new configuration.")
var NoCurrentContext = errors.New("No current context is set; cannot continue.\nSee 'help auth' for assistance creating/deleting/modifying/setting contexts.")
