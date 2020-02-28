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
	"fmt"
	"strings"

	apiTypes "github.com/liquidweb/liquidweb-cli/types/api"
)

// it helped package scope and getting access to the global client variable to
// put functions within this package... I'm not sure if a bigger package is
// more appropos that Client is a member of, or if this is actually fine..
// feedback is welcome

func (ci *Client) DerivePrivateParentUniqId(name string) (string, error) {
	var (
		privateParentUniqId     string
		privateParentDetails    apiTypes.CloudPrivateParentDetails
		privateParentDetailsErr error
	)

	// if name looks like a uniq_id, try it as a uniq_id first.
	if len(name) == 6 && strings.ToUpper(name) == name {
		if err := ci.CallLwApiInto("bleed/storm/private/parent/details",
			map[string]interface{}{"uniq_id": name},
			&privateParentDetails); err == nil {
			privateParentUniqId = name
		} else {
			privateParentDetailsErr = fmt.Errorf(
				"failed fetching parent details treating given --private-parent arg as a uniq-id [%s]: %s",
				name, err)
		}
	}

	// if we havent found the pp details yet, try assuming name is the name of the pp
	if privateParentUniqId == "" {
		methodArgs := AllPaginatedResultsArgs{
			Method:         "bleed/storm/private/parent/list",
			ResultsPerPage: 100,
		}
		results, err := ci.AllPaginatedResults(&methodArgs)
		if err != nil {
			ci.Die(err)
		}

		for _, item := range results.Items {
			var privateParentDetails apiTypes.CloudPrivateParentDetails
			if err := CastFieldTypes(item, &privateParentDetails); err != nil {
				ci.Die(err)
			}

			if privateParentDetails.Domain == name {
				// found it get details
				err := ci.CallLwApiInto("bleed/storm/private/parent/details",
					map[string]interface{}{
						"uniq_id": privateParentDetails.UniqId,
					},
					&privateParentDetails)
				if err != nil {
					privateParentDetailsErr = fmt.Errorf(
						"failed fetching private parent details for discovered uniq-id [%s] error: %s %w",
						privateParentDetails.UniqId, err, privateParentDetailsErr)
					return "", privateParentDetailsErr
				}
				privateParentUniqId = privateParentDetails.UniqId
				break // found the uniq_id so break
			}
		}
	}

	if privateParentUniqId == "" {
		return "", fmt.Errorf("failed deriving uniq_id of private parent from [%s]: %s", name, privateParentDetailsErr)
	}

	return privateParentUniqId, nil
}
