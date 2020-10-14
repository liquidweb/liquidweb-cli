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
	"errors"
	"fmt"

	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/validate"
)

type CloudServerResizeParams struct {
	UniqId        string `yaml:"uniq-id"`
	ConfigId      int64  `yaml:"config-id"`
	SkipFsResize  bool   `yaml:"skip-fs-resize"`
	PrivateParent string `yaml:"private-parent"`
	Memory        int64  `yaml:"memory"`
	Vcpu          int64  `yaml:"vcpu"`
	DiskSpace     int64  `yaml:"disk-space"`
}

func (self *CloudServerResizeParams) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// define defaults
	type rawType CloudServerResizeParams
	raw := rawType{
		ConfigId:  -1,
		Memory:    -1,
		Vcpu:      -1,
		DiskSpace: -1,
	} // Put your defaults here
	if err := unmarshal(&raw); err != nil {
		return err
	}
	*self = CloudServerResizeParams(raw)

	return nil
}

func (self *Client) CloudServerResize(params *CloudServerResizeParams) (result string, err error) {

	validateFields := map[interface{}]interface{}{
		params.UniqId: "UniqId",
	}

	// must validate UniqId now because we call api methods with this uniq_id before below validate
	if err = validate.Validate(validateFields); err != nil {
		return
	}

	// convert bool to int for api
	skipFsResizeInt := 0
	if params.SkipFsResize {
		skipFsResizeInt = 1
	}

	if params.ConfigId == -1 && params.PrivateParent == "" {
		err = errors.New("flag --config-id required when --private-parent is not given")
		return
	}

	resizeArgs := map[string]interface{}{
		"uniq_id":        params.UniqId,
		"skip_fs_resize": skipFsResizeInt,
		"newsize":        params.ConfigId,
	}

	// get details of existing configuration
	var cloudServerDetails apiTypes.CloudServerDetails
	if err = self.CallLwApiInto(
		"bleed/storm/server/details",
		map[string]interface{}{
			"uniq_id": params.UniqId,
		}, &cloudServerDetails); err != nil {
		return
	}

	var (
		liveResize      bool
		twoRebootResize bool
	)
	if params.PrivateParent == "" {
		// non private parent resize
		if params.Memory != -1 || params.DiskSpace != -1 || params.Vcpu != -1 {
			err = errors.New("cannot pass --memory --diskspace or --vcpu when --private-parent is not given")
			return
		}

		// if already on the given config, nothing to do
		if cloudServerDetails.ConfigId == params.ConfigId {
			err = fmt.Errorf("already on config-id [%d]; not initiating a resize", params.ConfigId)
			return
		}

		validateFields[params.ConfigId] = "PositiveInt64"
		if err = validate.Validate(validateFields); err != nil {
			return
		}

		// determine reboot expectation.
		//   resize up full: 2 reboot
		//   resize up quick (skip-fs-resize) 1 reboot
		//   resize down: 1 reboot
		var configDetails apiTypes.CloudConfigDetails
		if err = self.CallLwApiInto("bleed/storm/config/details",
			map[string]interface{}{"id": params.ConfigId}, &configDetails); err != nil {
			return
		}

		if configDetails.Disk >= cloudServerDetails.DiskSpace {
			// disk space going up..
			if !params.SkipFsResize {
				// .. and not skipping fs resize, will be 2 reboots.
				twoRebootResize = true
			}
		}
	} else {
		// private parent resize specific logic
		if params.Memory == -1 && params.DiskSpace == -1 && params.Vcpu == -1 {
			err = errors.New("resizes on private parents require at least least one of: --memory --diskspace --vcpu flags")
			return
		}

		var privateParentUniqId string
		privateParentUniqId, _, err = self.DerivePrivateParentUniqId(params.PrivateParent)
		if err != nil {
			return
		}

		var (
			diskspaceChanging bool
			vcpuChanging      bool
			memoryChanging    bool
			memoryCanLive     bool
			vcpuCanLive       bool
		)
		// record what resources are changing
		if params.DiskSpace != -1 {
			if cloudServerDetails.DiskSpace != params.DiskSpace {
				diskspaceChanging = true
			}
		}
		if params.Vcpu != -1 {
			if cloudServerDetails.Vcpu != params.Vcpu {
				vcpuChanging = true
			}
		}
		if params.Memory != -1 {
			if cloudServerDetails.Memory != params.Memory {
				memoryChanging = true
			}
		}
		// allow resizes to a private parent even if its old non private parent config had exact same specs
		if cloudServerDetails.ConfigId == 0 && cloudServerDetails.PrivateParent != privateParentUniqId {
			if !diskspaceChanging && !vcpuChanging && !memoryChanging {
				err = errors.New("private parent resize, but passed diskspace, memory, vcpu values match existing values")
				return
			}
		}

		resizeArgs["newsize"] = 0                  // 0 indicates private parent resize
		resizeArgs["parent"] = privateParentUniqId // uniq_id of the private parent
		validateFields[privateParentUniqId] = "UniqId"
		// server/resize api method always wants diskspace, vcpu, memory passed for pp resize, even if not changing
		// value. So set to current value, then override based on passed flags.
		resizeArgs["diskspace"] = cloudServerDetails.DiskSpace
		resizeArgs["memory"] = cloudServerDetails.Memory
		resizeArgs["vcpu"] = cloudServerDetails.Vcpu

		if params.DiskSpace != -1 {
			resizeArgs["diskspace"] = params.DiskSpace // desired diskspace
			validateFields[params.DiskSpace] = "PositiveInt64"
		}
		if params.Memory != -1 {
			resizeArgs["memory"] = params.Memory // desired memory
			validateFields[params.Memory] = "PositiveInt64"
		}
		if params.Vcpu != -1 {
			resizeArgs["vcpu"] = params.Vcpu // desired vcpus
			validateFields[params.Vcpu] = "PositiveInt64"
		}

		// determine if this will be a live resize
		if _, exists := resizeArgs["memory"]; exists {
			if params.Memory >= cloudServerDetails.Memory {
				// asking for more RAM
				memoryCanLive = true
			}
		}
		if _, exists := resizeArgs["vcpu"]; exists {
			if params.Vcpu >= cloudServerDetails.Vcpu {
				// asking for more vcpu
				vcpuCanLive = true
			}
		}

		if params.Memory != -1 && params.Vcpu != -1 {
			if vcpuCanLive && memoryCanLive {
				liveResize = true
			}
		} else if memoryCanLive {
			liveResize = true
		} else if vcpuCanLive {
			liveResize = true
		}

		// if diskspace allocation changes its not currently ever done live regardless of memory, vcpu
		if params.DiskSpace != -1 {
			if resizeArgs["diskspace"] != cloudServerDetails.DiskSpace {
				liveResize = false
			}
		}
	}

	if err = validate.Validate(validateFields); err != nil {
		return
	}

	if _, err = self.LwCliApiClient.Call("bleed/server/resize", resizeArgs); err != nil {
		return
	}

	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("server resized started! You can check progress with 'cloud server status --uniq-id %s'\n\n", params.UniqId))

	if liveResize {
		b.WriteString(fmt.Sprintf("\nthis resize will be performed live without downtime.\n"))
	} else {
		rebootExpectation := "one reboot"
		if twoRebootResize {
			rebootExpectation = "two reboots"
		}
		b.WriteString(fmt.Sprintf(
			"\nexpect %s during this process. Your server will be online as the disk is copied to the destination.\n",
			rebootExpectation))

		if twoRebootResize {
			b.WriteString(fmt.Sprintf(
				"\tTIP: Avoid the second reboot by passing --skip-fs-resize. See usage for additional details.\n"))
		}
	}

	result = b.String()

	return
}
