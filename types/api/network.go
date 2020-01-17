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
package apiTypes

import (
	"fmt"
	"strings"
)

type NetworkIpPoolListEntry struct {
	Id     int64 `json:"id" mapstructure:"id"`
	ZoneId int64 `json:"zone_id" mapstructure:"zone_id"`
}

type NetworkIpPoolDetails struct {
	Accnt       int64                            `json:"accnt" mapstructure:"accnt"`
	Id          int64                            `json:"id" mapstructure:"id"`
	UniqId      string                           `json:"uniq_id" mapstructure:"uniq_id"`
	ZoneId      int64                            `json:"zone_id" mapstructure:"zone_id"`
	Assignments []NetworkIpPoolDetailsAssignment `json:"assignments" mapstructure:"assignments"`
}

type NetworkIpPoolDetailsAssignment struct {
	BeginRange string `json:"begin_range" mapstructure:"begin_range"`
	Broadcast  string `json:"broadcast" mapstructure:"broadcast"`
	EndRange   string `json:"end_range" mapstructure:"end_range"`
	Gateway    string `json:"gateway" mapstructure:"gateway"`
	Id         int64  `json:"id" mapstructure:"id"`
	Netmask    string `json:"netmask" mapstructure:"netmask"`
	Network    string `json:"network" mapstructure:"network"`
	ZoneId     int64  `json:"zone_id" mapstructure:"zone_id"`
}

func (x NetworkIpPoolDetails) String() string {
	var slice []string

	slice = append(slice, fmt.Sprintf("IP Pool id [%d] uniq_id [%s]\n", x.Id, x.UniqId))
	slice = append(slice, fmt.Sprintf("\tZoneId: %d\n", x.ZoneId))
	slice = append(slice, fmt.Sprintf("\tAccount: %d\n", x.Accnt))
	slice = append(slice, fmt.Sprintf("\tAssignments:\n"))
	for _, assignment := range x.Assignments {
		slice = append(slice, fmt.Sprintf("\t\tassignment:\n"))
		slice = append(slice, fmt.Sprintf("\t\t\tBeginRange: %s\n", assignment.BeginRange))
		slice = append(slice, fmt.Sprintf("\t\t\tEndRange: %s\n", assignment.EndRange))
		if assignment.Broadcast != "" {
			slice = append(slice, fmt.Sprintf("\t\t\tBroadcast: %s\n", assignment.Broadcast))
		}
		slice = append(slice, fmt.Sprintf("\t\t\tGateway: %s\n", assignment.Gateway))
		slice = append(slice, fmt.Sprintf("\t\t\tNetmask: %s\n", assignment.Netmask))
		slice = append(slice, fmt.Sprintf("\t\t\tNetwork: %s\n", assignment.Network))
		slice = append(slice, fmt.Sprintf("\t\t\tId: %d\n", assignment.Id))
		slice = append(slice, fmt.Sprintf("\t\t\tZoneId: %d\n", assignment.ZoneId))
	}

	return strings.Join(slice[:], "")
}

type NetworkIpPoolDelete struct {
	Deleted bool `json:"deleted" mapstructure:"deleted"`
}

type NetworkIpAdd struct {
	Adding string `json:"adding" mapstructure:"adding"`
}

type NetworkIpRemove struct {
	Removing string `json:"removing" mapstructure:"removing"`
}

type NetworkAssignmentListEntry struct {
	Broadcast string `json:"broadcast" mapstructure:"broadcast"`
	Ip        string `json:"ip" mapstructure:"ip"`
	Gateway   string `json:"gateway" mapstructure:"gateway"`
	Id        int64  `json:"id" mapstructure:"id"`
	Netmask   string `json:"netmask" mapstructure:"netmask"`
	Network   string `json:"network" mapstructure:"network"`
}

func (x NetworkAssignmentListEntry) String() string {
	var slice []string

	slice = append(slice, fmt.Sprintf("\tIP: %s\n", x.Ip))
	slice = append(slice, fmt.Sprintf("\t\tId: %d\n", x.Id))
	slice = append(slice, fmt.Sprintf("\t\tGateway: %s\n", x.Gateway))
	slice = append(slice, fmt.Sprintf("\t\tBroadcast: %s\n", x.Broadcast))
	slice = append(slice, fmt.Sprintf("\t\tNetmask: %s\n", x.Netmask))
	slice = append(slice, fmt.Sprintf("\t\tNetwork: %s\n", x.Netmask))

	return strings.Join(slice[:], "")
}
