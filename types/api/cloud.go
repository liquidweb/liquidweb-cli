package apiTypes

type CloudServerStatus struct {
	DetailedStatus string                         `json: "detailed_status" mapstructure:"detailed_status"`
	Progress       float64                        `json: "progress" mapstructure:"progress"`
	Running        []CloudServerStatusRunningData `json: "running" mapstructure:"running"`
	Status         string                         `json:"status" mapstructure:"status"`
}

type CloudServerStatusRunningData struct {
	CurrentStep    int64  `json: "current_step" mapstructure: "current_step"`
	DetailedStatus string `json: "detailed_status" mapstructure: "detailed_status"`
	Name           string `json: "name" mapstructure: "name"`
	Status         string `json: "status" mapstructure: "status"`
}

type CloudServerRebootResponse struct {
	Rebooted string `json: "rebooted" mapstructure: "rebooted"`
}
