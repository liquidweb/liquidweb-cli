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

	"github.com/spf13/cast"

	apiTypes "github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/validate"
)

type CloudServerCreateParams struct {
	Template        string
	Type            string
	Hostname        string
	Ips             int64
	PoolIps         []string
	PublicSshKey    string
	ConfigId        int64
	BackupPlan      string
	BackupPlanQuota int64
	Bandwidth       string
	Zone            int64
	WinAv           string
	MsSql           string // windows
	PrivateParent   string
	Password        string
	Memory          int64 // required only if private parent
	Diskspace       int64 // required only if private parent
	Vcpu            int   // required only if private parent
	BackupId        int64 // create from backup
	ImageId         int64 // create from image
}

func (ci *Client) CloudServerCreate(params *CloudServerCreateParams) string {
	var err error

	// sanity check flags
	if params.PrivateParent != "" {
		if params.ConfigId > 0 {
			ci.Die(fmt.Errorf("--config_id must be 0 or omitted when specifying --private-parent"))
		}
		// create on a private parent. diskspace, memory, vcpu are required.
		if params.Memory == -1 {
			ci.Die(fmt.Errorf("--memory is required when specifying --private-parent"))
		}
		if params.Diskspace == -1 {
			ci.Die(fmt.Errorf("--diskspace is required when specifying --private-parent"))
		}
		if params.Vcpu == -1 {
			ci.Die(fmt.Errorf("--vcpu is required when specifying --private-parent"))
		}
	} else {
		if params.ConfigId <= 0 {
			ci.Die(fmt.Errorf("--config_id is required when not specifying --private-parent"))
		}

	}

	if params.Template == "" && params.BackupId == -1 && params.ImageId == -1 {
		ci.Die(fmt.Errorf("at least one of the following flags must be set --template --image-id --backup-id"))
	}

	// TODO - not sure if input validation belongs here or in the command...

	validateFields := map[interface{}]interface{}{
		params.Zone:       map[string]string{"type": "PositiveInt", "optional": "true"},
		params.Hostname:   "NonEmptyString",
		params.Type:       "NonEmptyString",
		params.Ips:        "PositiveInt",
		params.Password:   "NonEmptyString",
		params.BackupPlan: "NonEmptyString",
	}
	if params.BackupId != -1 {
		validateFields[params.BackupId] = "PositiveInt"
	}
	if params.ImageId != -1 {
		validateFields[params.ImageId] = "PositiveInt"
	}
	if params.Vcpu == -1 {
		validateFields[params.ConfigId] = "PositiveInt"
	}
	if params.ConfigId == -1 {
		validateFields[params.Vcpu] = "PositiveInt"
		validateFields[params.Memory] = "PositiveInt"
		validateFields[params.Diskspace] = "PositiveInt"
	}
	if err := validate.Validate(validateFields); err != nil {
		ci.Die(err)
	}

	// if passed a private-parent flag, derive its uniq_id
	var privateParentUniqId string
	if params.PrivateParent != "" {
		privateParentUniqId, err = ci.DerivePrivateParentUniqId(params.PrivateParent)
		if err != nil {
			ci.Die(err)
		}
	}

	// buildout args for bleed/server/create
	createArgs := map[string]interface{}{
		"domain":   params.Hostname,
		"type":     params.Type,
		"pool_ips": params.PoolIps,
		"new_ips":  params.Ips,
		"zone":     params.Zone,
		"password": params.Password,
		"features": map[string]interface{}{
			"Bandwidth": params.Bandwidth,
			"ConfigId":  params.ConfigId,
			"ExtraIp": map[string]interface{}{
				"value": params.Ips,
				"count": 0,
			},
			"LiquidWebBackupPlan": params.BackupPlan,
		},
	}

	var isWindows bool
	if params.Template != "" {
		createArgs["features"].(map[string]interface{})["Template"] = params.Template
		if strings.Contains(strings.ToUpper(params.Template), "WINDOWS") {
			isWindows = true
		}
	}
	if params.BackupId != -1 {
		// check backup and see if its windows
		apiArgs := map[string]interface{}{"id": params.BackupId}
		var details apiTypes.CloudBackupDetails
		err := ci.CallLwApiInto("bleed/storm/backup/details", apiArgs, &details)
		if err != nil {
			ci.Die(err)
		}
		if strings.Contains(strings.ToUpper(details.Template), "WINDOWS") {
			isWindows = true
		}
		createArgs["backup_id"] = params.BackupId
	}
	if params.ImageId != -1 {
		// check image and see if its windows
		apiArgs := map[string]interface{}{"id": params.ImageId}
		var details apiTypes.CloudImageDetails
		err := ci.CallLwApiInto("bleed/storm/image/details", apiArgs, &details)
		if err != nil {
			ci.Die(err)
		}
		if strings.Contains(strings.ToUpper(details.Template), "WINDOWS") {
			isWindows = true
		}
		createArgs["image_id"] = params.ImageId
	}

	// windows servers need special arguments
	if isWindows {
		if params.WinAv == "" {
			params.WinAv = "None"
		}
		createArgs["features"].(map[string]interface{})["WinAV"] = params.WinAv
		createArgs["features"].(map[string]interface{})["WindowsLicense"] = "Windows"
		if params.Type == "SS.VPS" {
			createArgs["type"] = "SS.VPS.WIN"
		}
		if params.MsSql == "" {
			params.MsSql = "None"
		}
		var coreCnt int
		if params.Vcpu == -1 {
			// standard config_id create, fetch configs core count and use it
			var details apiTypes.CloudConfigDetails
			if err := ci.CallLwApiInto("bleed/storm/config/details",
				map[string]interface{}{"id": params.ConfigId}, &details); err != nil {
				ci.Die(err)
			}
			coreCnt = cast.ToInt(details.Vcpu)
		} else {
			// private parent, use vcpu flag
			coreCnt = params.Vcpu
		}
		createArgs["features"].(map[string]interface{})["MsSQL"] = map[string]interface{}{
			"value": params.MsSql,
			"count": coreCnt,
		}
	}

	if privateParentUniqId != "" {
		createArgs["parent"] = privateParentUniqId
		createArgs["vcpu"] = params.Vcpu
		createArgs["diskspace"] = params.Diskspace
		createArgs["memory"] = params.Memory
	}

	if params.BackupPlan == "Quota" {
		createArgs["features"].(map[string]interface{})["BackupQuota"] = params.BackupPlanQuota
	}

	if params.PublicSshKey != "" {
		createArgs["public_ssh_key"] = params.PublicSshKey
	}

	result, err := ci.LwCliApiClient.Call("bleed/server/create", createArgs)
	if err != nil {
		ci.Die(err)
	}

	return cast.ToString(result.(map[string]interface{})["uniq_id"])
}
