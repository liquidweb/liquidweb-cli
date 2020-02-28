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
	"github.com/liquidweb/liquidweb-cli/utils"
	"github.com/liquidweb/liquidweb-cli/validate"
)

type CloudServerCreateParams struct {
	Template        string   `yaml:"template"`
	Type            string   `yaml:"type"`
	Hostname        string   `yaml:"hostname"`
	Ips             int      `yaml:"ips"`
	PoolIps         []string `yaml:"pool-ips"`
	PublicSshKey    string   `yaml:"public-ssh-key"`
	ConfigId        int      `yaml:"config-id"`
	BackupPlan      string   `yaml:"backup-plan"`
	BackupPlanQuota int      `yaml:"backup-plan-quota"`
	Bandwidth       string   `yaml:"bandwidth"`
	Zone            int      `yaml:"zone"`
	WinAv           string   `yaml:"winav"`  // windows
	MsSql           string   `yaml:"ms-sql"` // windows
	PrivateParent   string   `yaml:"private-parent"`
	Password        string   `yaml:"password"`
	Memory          int      `yaml:"memory"`    // required only if private parent
	Diskspace       int      `yaml:"diskspace"` // required only if private parent
	Vcpu            int      `yaml:"vcpu"`      // required only if private parent
	BackupId        int      `yaml:"backup-id"` //create from backup
	ImageId         int      `yaml:"image-id"`  // create from image
}

func (s *CloudServerCreateParams) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// define defaults
	type rawType CloudServerCreateParams
	raw := rawType{
		BackupId:   -1,
		ImageId:    -1,
		Vcpu:       -1,
		Memory:     -1,
		Diskspace:  -1,
		Bandwidth:  "SS.5000",
		BackupPlan: "None",
		Ips:        1,
		Type:       "SS.VPS",
	} // Put your defaults here
	if err := unmarshal(&raw); err != nil {
		return err
	}
	*s = CloudServerCreateParams(raw)
	return nil
}

func (ci *Client) CloudServerCreate(params *CloudServerCreateParams) (string, error) {
	var err error

	// if passed a private-parent flag, derive its uniq_id
	if params.PrivateParent != "" {
		params.PrivateParent, err = ci.DerivePrivateParentUniqId(params.PrivateParent)
		if err != nil {
			return "", err
		}
	}

	// default password
	if params.Password == "" {
		params.Password = utils.RandomString(25)
	}

	//default hostname
	if params.Hostname == "" {
		params.Hostname = fmt.Sprintf("%s.%s.io", utils.RandomString(4), utils.RandomString(10))
	}

	// sanity check flags
	if params.PrivateParent != "" {
		if params.ConfigId > 0 {
			return "", fmt.Errorf("--config_id must be 0 or omitted when specifying --private-parent")
		}
		// create on a private parent. diskspace, memory, vcpu are required.
		if params.Memory == -1 {
			return "", fmt.Errorf("--memory is required when specifying --private-parent")
		}
		if params.Diskspace == -1 {
			return "", fmt.Errorf("--diskspace is required when specifying --private-parent")
		}
		if params.Vcpu == -1 {
			return "", fmt.Errorf("--vcpu is required when specifying --private-parent")
		}
	} else {
		if params.ConfigId <= 0 {
			return "", fmt.Errorf("--config_id is required when not specifying --private-parent")
		}

	}

	if params.Template == "" && params.BackupId == -1 && params.ImageId == -1 {
		return "", fmt.Errorf("at least one of the following flags must be set --template --image-id --backup-id")
	}

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
		return "", err
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
			return "", err
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
			return "", err
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
				return "", err
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

	if params.PrivateParent != "" {
		createArgs["parent"] = params.PrivateParent
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
		return "", err
	}

	return cast.ToString(result.(map[string]interface{})["uniq_id"]), nil
}
