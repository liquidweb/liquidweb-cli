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

	resizePlanArgs := map[string]interface{}{
		"uniq_id": params.UniqId,
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

	var privateParentUniqId string

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

		resizePlanArgs["config_id"] = params.ConfigId
	} else {
		// private parent resize specific logic
		if params.Memory == -1 && params.DiskSpace == -1 && params.Vcpu == -1 {
			err = errors.New("resizes on private parents require at least least one of: --memory --diskspace --vcpu flags")
			return
		}

		privateParentUniqId, err = self.DerivePrivateParentUniqId(params.PrivateParent)
		if err != nil {
			return
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

		resizePlanArgs["config_id"] = 0
		resizePlanArgs["private_parent"] = privateParentUniqId
		resizePlanArgs["memory"] = resizeArgs["memory"]
		resizePlanArgs["disk"] = resizeArgs["diskspace"]
		resizePlanArgs["vcpu"] = resizeArgs["vcpu"]
	}

	if err = validate.Validate(validateFields); err != nil {
		return
	}

	var expectation apiTypes.CloudServerResizeExpectation
	if err = self.CallLwApiInto("bleed/storm/server/resizePlan", resizePlanArgs, &expectation); err != nil {
		err = fmt.Errorf("Configuration Not Available\n\n%s\n", err)
		return
	}

	if _, err = self.LwCliApiClient.Call("bleed/server/resize", resizeArgs); err != nil {
		return
	}

	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("Server resized started! You can check progress with 'cloud server status --uniq-id %s'\n\n", params.UniqId))
	b.WriteString(fmt.Sprintf("Resource changes: Memory [%d] Disk [%d] Vcpu [%d]\n", expectation.MemoryDifference,
		expectation.DiskDifference, expectation.VcpuDifference))

	if expectation.RebootRequired {
		b.WriteString("\nExpect a reboot during this resize.\n")
	} else {
		b.WriteString("\nThis resize will be performed live without downtime.\n")
	}

	result = b.String()

	return
}
