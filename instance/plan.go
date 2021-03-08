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
	Ssh   []SshParams
}

type PlanCloud struct {
	Server   *PlanCloudServer
	Template *PlanCloudTemplate
	Network  *PlanCloudNetwork
}

type PlanCloudServer struct {
	Create []CloudServerCreateParams
	Resize []CloudServerResizeParams
	Reboot []CloudServerRebootParams
}

type PlanCloudTemplate struct {
	Restore []CloudTemplateRestoreParams
}

type PlanCloudNetwork struct {
	Public  *PlanCloudNetworkPublic
	Private *PlanCloudNetworkPrivate
}

type PlanCloudNetworkPublic struct {
	Add    []CloudNetworkPublicAddParams
	Remove []CloudNetworkPublicRemoveParams
}

type PlanCloudNetworkPrivate struct {
	Attach []CloudNetworkPrivateAttachParams
	Detach []CloudNetworkPrivateDetachParams
}

func (ci *Client) ProcessPlan(plan *Plan) error {

	if plan.Cloud != nil {
		if err := ci.processPlanCloud(plan.Cloud); err != nil {
			return err
		}
	}

	for i, _ := range plan.Ssh {
		if err := ci.processPlanSsh(&plan.Ssh[i]); err != nil {
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

	if cloud.Network != nil {
		if err := ci.processPlanCloudNetwork(cloud.Network); err != nil {
			return err
		}
	}

	return nil
}

func (ci *Client) processPlanSsh(params *SshParams) (err error) {
	err = ci.Ssh(params)

	return
}

func (ci *Client) processPlanCloudServer(server *PlanCloudServer) error {

	if server.Create != nil {
		for i, _ := range server.Create {
			if err := ci.processPlanCloudServerCreate(&server.Create[i]); err != nil {
				return err
			}
		}
	}
	if server.Resize != nil {
		for i, _ := range server.Resize {
			if err := ci.processPlanCloudServerResize(&server.Resize[i]); err != nil {
				return err
			}
		}
	}
	if server.Reboot != nil {
		for i, _ := range server.Reboot {
			if err := ci.processPlanCloudServerReboot(&server.Reboot[i]); err != nil {
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

func (ci *Client) processPlanCloudServerReboot(params *CloudServerRebootParams) error {

	result, err := ci.CloudServerReboot(params)
	if err != nil {
		return err
	}

	fmt.Print(result)

	return nil
}

func (ci *Client) processPlanCloudTemplate(template *PlanCloudTemplate) error {

	if template.Restore != nil {
		for i, _ := range template.Restore {
			if err := ci.processPlanCloudTemplateRestore(&template.Restore[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ci *Client) processPlanCloudNetwork(network *PlanCloudNetwork) error {

	if network.Public != nil {
		if err := ci.processPlanCloudNetworkPublic(network.Public); err != nil {
			return err
		}
	}

	if network.Private != nil {
		if err := ci.processPlanCloudNetworkPrivate(network.Private); err != nil {
			return err
		}
	}

	return nil
}

func (ci *Client) processPlanCloudNetworkPublic(public *PlanCloudNetworkPublic) error {

	if public.Add != nil {
		for i, _ := range public.Add {
			if err := ci.processPlanCloudNetworkPublicAdd(&public.Add[i]); err != nil {
				return err
			}
		}
	}

	if public.Remove != nil {
		for i, _ := range public.Remove {
			if err := ci.processPlanCloudNetworkPublicRemove(&public.Remove[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ci *Client) processPlanCloudNetworkPrivate(private *PlanCloudNetworkPrivate) error {

	if private.Attach != nil {
		for i, _ := range private.Attach {
			if err := ci.processPlanCloudNetworkPrivateAttach(&private.Attach[i]); err != nil {
				return err
			}
		}
	}

	if private.Detach != nil {
		for i, _ := range private.Detach {
			if err := ci.processPlanCloudNetworkPrivateDetach(&private.Detach[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ci *Client) processPlanCloudNetworkPublicAdd(params *CloudNetworkPublicAddParams) error {
	result, err := ci.CloudNetworkPublicAdd(params)
	if err != nil {
		ci.Die(err)
	}

	fmt.Print(result)

	return nil
}

func (ci *Client) processPlanCloudNetworkPublicRemove(params *CloudNetworkPublicRemoveParams) error {
	result, err := ci.CloudNetworkPublicRemove(params)
	if err != nil {
		ci.Die(err)
	}

	fmt.Print(result)

	return nil
}

func (ci *Client) processPlanCloudNetworkPrivateAttach(params *CloudNetworkPrivateAttachParams) error {
	result, err := ci.CloudNetworkPrivateAttach(params)
	if err != nil {
		ci.Die(err)
	}

	fmt.Print(result)

	return nil
}

func (ci *Client) processPlanCloudNetworkPrivateDetach(params *CloudNetworkPrivateDetachParams) error {
	result, err := ci.CloudNetworkPrivateDetach(params)
	if err != nil {
		ci.Die(err)
	}

	fmt.Print(result)

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
