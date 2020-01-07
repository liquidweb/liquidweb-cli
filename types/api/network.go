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
	Primary    bool   `json:"primary" mapstructure:"primary"`
	ZoneId     int64  `json:"zone_id" mapstructure:"zone_id"`
}

type NetworkIpPoolDelete struct {
	Deleted bool `json:"deleted" mapstructure:"deleted"`
}
