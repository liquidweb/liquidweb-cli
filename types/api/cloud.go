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

type CloudServerStatus struct {
	DetailedStatus string                         `json:"detailed_status" mapstructure:"detailed_status"`
	Progress       float64                        `json:"progress" mapstructure:"progress"`
	Running        []CloudServerStatusRunningData `json:"running" mapstructure:"running"`
	Status         string                         `json:"status" mapstructure:"status"`
}

type CloudServerStatusRunningData struct {
	CurrentStep    int64  `json:"current_step" mapstructure:"current_step"`
	DetailedStatus string `json:"detailed_status" mapstructure:"detailed_status"`
	Name           string `json:"name" mapstructure:"name"`
	Status         string `json:"status" mapstructure:"status"`
}

type CloudServerRebootResponse struct {
	Rebooted string `json:"rebooted" mapstructure:"rebooted"`
}

type CloudServerDetails struct {
	Accnt               int64                  `json:"accnt" mapstructure:"accnt"`
	ConfigId            int64                  `json:"config_id" mapstructure:"config_id"`
	Memory              int64                  `json:"memory" mapstructure:"memory"`
	Template            string                 `json:"template" mapstructure:"template"`
	Type                string                 `json:"type" mapstructure:"type"`
	BackupEnabled       int64                  `json:"backup_enabled" mapstructure:"backup_enabled"`
	BackupSize          float64                `json:"backup_size" mapstructure:"backup_size"`
	UniqId              string                 `json:"uniq_id" mapstructure:"uniq_id"`
	Vcpu                int64                  `json:"vcpu" mapstructure:"vcpu"`
	BackupPlan          string                 `json:"backup_plan" mapstructure:"backup_plan"`
	BandwidthQuota      string                 `json:"bandwidth_quota" mapstructure:"bandwidth_quota"`
	Ip                  string                 `json:"ip" mapstructure:"ip"`
	IpCount             int64                  `json:"ip_count" mapstructure:"ip_count"`
	ManageLevel         string                 `json:"manage_level" mapstructure:"manage_level"`
	CreateDate          string                 `json:"create_date" mapstructure:"create_date"`
	DiskSpace           int64                  `json:"diskspace" mapstructure:"diskspace"`
	Domain              string                 `json:"domain" mapstructure:"domain"`
	Active              int64                  `json:"active" mapstructure:"active"`
	BackupQuota         int64                  `json:"backup_quota" mapstructure:"backup_quota"`
	Zone                CloudServerDetailsZone `json:"zone" mapstructure:"zone"`
	ConfigDescription   string                 `json:"config_description" mapstructure:"config_description"`
	TemplateDescription string                 `json:"template_description" mapstructure:"template_description"`
}

type CloudServerDetailsZone struct {
	Id     int64                        `json:"id" mapstructure:"id"`
	Name   string                       `json:"name" mapstructure:"name"`
	Region CloudServerDetailsZoneRegion `json:"region" mapstructure:"region"`
}

type CloudServerDetailsZoneRegion struct {
	Id   int64  `json:"id" mapstructure:"id"`
	Name string `json:"name" mapstructure:"name"`
}

type CloudPrivateParentDetails struct {
	Accnt             int64                                     `json:"accnt" mapstructure:"accnt"`
	BucketUniqId      string                                    `json:"bucket_uniq_id" mapstructure:"bucket_uniq_id"`
	ConfigDescription string                                    `json:"config_description" mapstructure:"config_description"`
	ConfigId          int64                                     `json:"config_id" mapstructure:"config_id"`
	CreateDate        string                                    `json:"create_date" mapstructure:"create_date"`
	DiskDetails       CloudPrivateParentDetailsEntryDiskDetails `json:"diskDetails" mapstructure:"diskDetails"`
	Domain            string                                    `json:"domain" mapstructure:"domain"`
	Id                int64                                     `json:"id" mapstructure:"id"`
	LicenseState      string                                    `json:"license_state" mapstructure:"license_state"`
	RegionId          int64                                     `json:"region_id" mapstructure:"region_id"`
	Resources         CloudPrivateParentDetailsEntryResource    `json:"resources" mapstructure:"resources"`
	SalesforceAsset   string                                    `json:"salesforce_asset" mapstructure:"salesforce_asset"`
	Status            string                                    `json:"status" mapstructure:"status"`
	Subaccnt          int64                                     `json:"subaccnt" mapstructure:"subaccnt"`
	Type              string                                    `json:"type" mapstructure:"type"`
	UniqId            string                                    `json:"uniq_id" mapstructure:"uniq_id"`
	Vcpu              int64                                     `json:"vcpu" mapstructure:"vcpu"`
	Zone              CloudPrivateParentDetailsEntryZone        `json:"zone" mapstructure:"zone"`
}

type CloudPrivateParentDetailsEntryResource struct {
	DiskSpace CloudPrivateParentDetailsEntryResourceEntry `json:"diskspace" mapstructure:"diskspace"`
	Memory    CloudPrivateParentDetailsEntryResourceEntry `json:"memory" mapstructure:"memory"`
}

type CloudPrivateParentDetailsEntryResourceEntry struct {
	Free  int64 `json:"free" mapstructure:"free"`
	Total int64 `json:"total" mapstructure:"total"`
	Used  int64 `json:"used" mapstructure:"used"`
}

type CloudPrivateParentDetailsEntryDiskDetails struct {
	Allocated int64 `json:"allocated" mapstructure:"allocated"`
	Snapshots int64 `json:"snapshots" mapstructure:"snapshots"`
}

type CloudPrivateParentDetailsEntryZone struct {
	AvailabilityZone string                                   `json:"availability_zone" mapstructure:"availability_zone"`
	Description      string                                   `json:"description" mapstructure:"description"`
	HvType           string                                   `json:"hv_type" mapstructure:"hv_type"`
	Id               int64                                    `json:"id" mapstructure:"id"`
	Legacy           int64                                    `json:"legacy" mapstructure:"legacy"`
	Name             string                                   `json:"name" mapstructure:"name"`
	Region           CloudPrivateParentDetailsEntryZoneRegion `json:"region" mapstructure:"region"`
	Status           string                                   `json:"status" mapstructure:"status"`
	ValidSourceHvs   []string                                 `json:"valid_source_hvs" mapstructure:"valid_source_hvs"`
}

type CloudPrivateParentDetailsEntryZoneRegion struct {
	Id   int64  `json:"id" mapstructure:"id"`
	Name string `json:"name" mapstructure:"name"`
}

type CloudConfigDetails struct {
	Id                int64         `json:"id" mapstructure:"id"`
	Active            int64         `json:"active" mapstructure:"active"`
	Available         int64         `json:"available" mapstructure:"available"`
	Category          string        `json:"category" mapstructure:"category"`
	Description       string        `json:"description" mapstructure:"description"`
	Disk              int64         `json:"disk,omitempty" mapstructure:"disk"`
	Featured          int64         `json:"featured" mapstructure:"featured"`
	Memory            int64         `json:"memory,omitempty" mapstructure:"memory"`
	Vcpu              int64         `json:"vcpu,omitempty" mapstructure:"vcpu"`
	ZoneAvailability  []map[int]int `json:"zone_availability" mapstructure:"zone_availability"`
	Retired           int64         `json:"retired,omitempty" mapstructure:"retired"`
	RamTotal          int64         `json:"ram_total,omitempty" mapstructure:"ram_total"`
	RamAvailable      int64         `json:"ram_available,omitempty" mapstructure:"ram_available"`
	RaidLevel         int64         `json:"raid_level,omitempty" mapstructure:"raid_level"`
	DiskType          int64         `json:"disk_type,omitempty" mapstructure:"disk_type"`
	DiskTotal         int64         `json:"disk_total,omitempty" mapstructure:"disk_total"`
	DiskCount         int64         `json:"disk_count,omitempty" mapstructure:"disk_count"`
	CpuSpeed          int64         `json:"cpu_speed,omitempty" mapstructure:"cpu_speed"`
	CpuModel          int64         `json:"cpu_model,omitempty" mapstructure:"cpu_model"`
	CpuHyperthreading int64         `json:"cpu_hyperthreading,omitempty" mapstructure:"cpu_hyperthreading"`
	CpuCount          int64         `json:"cpu_count,omitempty" mapstructure:"cpu_count"`
	CpuCores          int64         `json:"cpu_cores,omitempty" mapstructure:"cpu_cores"`
}

type CloudServerDestroyResponse struct {
	Destroyed string `json:"destroyed" mapstructure:"destroyed"`
}

type CloudServerShutdownResponse struct {
	Shutdown string `json:"shutdown" mapstructure:"shutdown"`
}

type CloudServerStartResponse struct {
	Started string `json:"started" mapstructure:"started"`
}
