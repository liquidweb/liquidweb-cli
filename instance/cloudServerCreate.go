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

	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/utils"
	"github.com/liquidweb/liquidweb-cli/validate"
)

type CloudServerCreateParams struct {
	Template      string   `yaml:"template"`
	Type          string   `yaml:"type"`
	Hostname      string   `yaml:"hostname"`
	Ips           int      `yaml:"ips"`
	PoolIps       []string `yaml:"pool-ips"`
	PublicSshKey  string   `yaml:"public-ssh-key"`
	ConfigId      int      `yaml:"config-id"`
	BackupDays    int      `yaml:"backup-days"`  // daily backup plan; how many days to keep a backup
	BackupQuota   int      `yaml:"backup-quota"` // backup quota plan; how many gb of backups to keep
	Bandwidth     string   `yaml:"bandwidth"`
	Zone          int      `yaml:"zone"`
	WinAv         string   `yaml:"winav"`  // windows
	MsSql         string   `yaml:"ms-sql"` // windows
	PrivateParent string   `yaml:"private-parent"`
	Password      string   `yaml:"password"`
	Memory        int      `yaml:"memory"`    // required only if private parent
	Diskspace     int      `yaml:"diskspace"` // required only if private parent
	Vcpu          int      `yaml:"vcpu"`      // required only if private parent
	BackupId      int      `yaml:"backup-id"` //create from backup
	ImageId       int      `yaml:"image-id"`  // create from image
}

func (s *CloudServerCreateParams) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// define defaults
	type rawType CloudServerCreateParams
	raw := rawType{
		BackupId:    -1,
		BackupDays:  -1,
		BackupQuota: -1,
		ImageId:     -1,
		Vcpu:        -1,
		Memory:      -1,
		Diskspace:   -1,
		Bandwidth:   "SS.5000",
		Ips:         1,
		Type:        "SS.VPS",
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
			return "", fmt.Errorf("--config-id must be 0 or omitted when specifying --private-parent")
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
			return "", fmt.Errorf("--config-id is required when not specifying --private-parent")
		}
		// not on a private parent, shouldnt pass private parent flags
		if params.Memory != -1 {
			return "", fmt.Errorf("--memory should not be passed with --config-id")
		}
		if params.Diskspace != -1 {
			return "", fmt.Errorf("--diskspace should not be passed with --config-id")
		}
		if params.Vcpu != -1 {
			return "", fmt.Errorf("--vcpu should not be passed with --config-id")
		}
	}

	if params.Template == "" && params.BackupId == -1 && params.ImageId == -1 {
		return "", fmt.Errorf("at least one of the following flags must be set --template --image-id --backup-id")
	}

	if params.BackupDays != -1 && params.BackupQuota != -1 {
		return "", fmt.Errorf("flags --backup-days and --backup-quota conflict")
	}

	validateFields := map[interface{}]interface{}{
		params.Zone:     map[string]string{"type": "PositiveInt", "optional": "true"},
		params.Hostname: "NonEmptyString",
		params.Type:     "NonEmptyString",
		params.Ips:      "PositiveInt",
		params.Password: "NonEmptyString",
	}
	if params.BackupId != -1 {
		validateFields[params.BackupId] = "PositiveInt"
	}
	if params.ImageId != -1 {
		validateFields[params.ImageId] = "PositiveInt"
	}
	if params.Vcpu == -1 {
		validateFields[params.Vcpu] = "PositiveInt"
	}
	if params.ConfigId == -1 {
		validateFields[params.Vcpu] = "PositiveInt"
		validateFields[params.Memory] = "PositiveInt"
		validateFields[params.Diskspace] = "PositiveInt"
	}
	if err := validate.Validate(validateFields); err != nil {
		return "", err
	}

	cloudBackupPlan := "None"
	if params.BackupDays != -1 {
		cloudBackupPlan = "Daily"
	} else if params.BackupQuota != -1 {
		cloudBackupPlan = "Quota"
	}

	// buildout args for bleed/server/create
	createArgs := map[string]interface{}{
		"domain":   params.Hostname,
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
			"LiquidWebBackupPlan": cloudBackupPlan,
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

	// when creating with a config-id, adjust the Type param for bare-metal types if blatantly wrong
	var configDetails apiTypes.CloudConfigDetails
	if params.ConfigId > 0 {
		if err := ci.CallLwApiInto("bleed/storm/config/details",
			map[string]interface{}{"id": params.ConfigId}, &configDetails); err != nil {
			return "", err
		}
		if configDetails.Category == "bare-metal" {
			if isWindows {
				params.Type = "SS.VM.WIN"
			} else {
				params.Type = "SS.VM"
			}
		} else if configDetails.Category == "bare-metal-r" {
			params.Type = "SS.VM.R"
		}
	}

	// windows servers need special arguments
	if isWindows {
		if params.WinAv == "" {
			params.WinAv = "None"
		}
		createArgs["features"].(map[string]interface{})["WinAV"] = params.WinAv
		createArgs["features"].(map[string]interface{})["WindowsLicense"] = "Windows"
		if params.Type == "SS.VPS" {
			params.Type = "SS.VPS.WIN"
		}
		if params.MsSql == "" {
			params.MsSql = "None"
		}
		var coreCnt int
		if params.Vcpu == -1 {
			// standard config_id create, fetch configs core count and use it
			coreCnt = cast.ToInt(configDetails.Vcpu)
		} else {
			// private parent, use vcpu flag
			coreCnt = params.Vcpu
		}
		createArgs["features"].(map[string]interface{})["MsSQL"] = map[string]interface{}{
			"value": params.MsSql,
			"count": coreCnt,
		}
	}

	createArgs["type"] = params.Type

	if params.PrivateParent != "" {
		createArgs["parent"] = params.PrivateParent
		createArgs["vcpu"] = params.Vcpu
		createArgs["diskspace"] = params.Diskspace
		createArgs["memory"] = params.Memory
	}

	if cloudBackupPlan == "Quota" {
		createArgs["features"].(map[string]interface{})["BackupQuota"] = params.BackupQuota
	} else if cloudBackupPlan == "Daily" {
		createArgs["features"].(map[string]interface{})["BackupDay"] = map[string]int{
			"value":     1,
			"num_units": params.BackupDays,
		}
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
