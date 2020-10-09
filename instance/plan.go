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
)

type Plan struct {
	Cloud *PlanCloud
}

type PlanCloud struct {
	Server   *PlanCloudServer
	Template *PlanCloudTemplate
}

type PlanCloudServer struct {
	Create []CloudServerCreateParams
	Resize []CloudServerResizeParams
}

type PlanCloudTemplate struct {
	Restore []CloudTemplateRestoreParams
}

func (ci *Client) ProcessPlan(plan *Plan) error {

	if plan.Cloud != nil {
		if err := ci.processPlanCloud(plan.Cloud); err != nil {
			return err
		}
	}

	return nil
}

func (ci *Client) processPlanCloud(cloud *PlanCloud) error {

	if cloud.Server != nil {
		if err := ci.processPlanCloudServer(cloud.Server); err != nil {
			return err
		}
	}

	if cloud.Template != nil {
		if err := ci.processPlanCloudTemplate(cloud.Template); err != nil {
			return err
		}
	}

	return nil
}

func (ci *Client) processPlanCloudServer(server *PlanCloudServer) error {

	if server.Create != nil {
		for _, c := range server.Create {
			if err := ci.processPlanCloudServerCreate(&c); err != nil {
				return err
			}
		}
	}
	if server.Resize != nil {
		for _, r := range server.Resize {
			if err := ci.processPlanCloudServerResize(&r); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ci *Client) processPlanCloudServerCreate(params *CloudServerCreateParams) error {

	uniqId, err := ci.CloudServerCreate(params)
	if err != nil {
		return err
	}

	fmt.Printf(
		"Cloud server with uniq-id [%s] creating. Check status with 'cloud server status --uniq-id %s'\n",
		uniqId, uniqId)
	return nil
}

func (ci *Client) processPlanCloudServerResize(params *CloudServerResizeParams) error {

	result, err := ci.CloudServerResize(params)
	if err != nil {
		return err
	}

	fmt.Print(result)

	return nil
}

func (ci *Client) processPlanCloudTemplate(template *PlanCloudTemplate) error {

	if template.Restore != nil {
		for _, c := range template.Restore {
			if err := ci.processPlanCloudTemplateRestore(&c); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ci *Client) processPlanCloudTemplateRestore(params *CloudTemplateRestoreParams) error {

	result, err := ci.CloudTemplateRestore(params)
	if err != nil {
		ci.Die(err)
	}

	fmt.Printf("Restoring template! %s\n", result)
	fmt.Printf("\tcheck progress with 'cloud server status --uniq-id %s'\n", params.UniqId)

	return nil
}
