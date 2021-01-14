package defaults

import (
	"fmt"
	"strings"
)

type AllFlags map[string]interface{}

func (self AllFlags) String() string {
	var slice []string

	if len(self) == 0 {
		slice = append(slice, "No configured default flags. Set some with 'default-flags set'.\n")
	} else {
		slice = append(slice, "Configured default flags:\n\n")

		for flag, value := range self {
			slice = append(slice, fmt.Sprintf("  Flag: %s\n", flag))
			slice = append(slice, fmt.Sprintf("    Value: %+v\n", value))
		}
	}

	return strings.Join(slice[:], "")
}

var permittedFlags = map[string]map[string]interface{}{
	// cloud network vip create
	"cloud_network_vip_create_zone": map[string]interface{}{
		"enabled":   true,
		"validator": "PositiveInt64",
	},
	// cloud private-parent create
	"cloud_private-parent_create_config-id": map[string]interface{}{
		"enabled": true,
		"type":    "PositiveInt64",
	},
	"cloud_private-parent_create_zone": map[string]interface{}{
		"enabled": true,
		"type":    "PositiveInt64",
	},
	// cloud server clone
	"cloud_server_clone_config-id": map[string]interface{}{
		"enabled": true,
		"type":    "PositiveInt64",
	},
	// cloud server create
	"cloud_server_create_zone": map[string]interface{}{
		"enabled": true,
		"type":    "PositiveInt64",
	},
	"cloud_server_create_template": map[string]interface{}{
		"enabled": true,
		"type":    "NonEmptyString",
	},
	"cloud_server_create_config-id": map[string]interface{}{
		"enabled": true,
		"type":    "PositiveInt64",
	},
	// cloud server resize
	"cloud_server_resize_config-id": map[string]interface{}{
		"enabled": true,
		"type":    "PositiveInt64",
	},
	// cloud template restore
	"cloud_template_restore_template": map[string]interface{}{
		"enabled": true,
		"type":    "NonEmptyString",
	},
	// network ip-pool create
	"network_ip-pool_create_zone": map[string]interface{}{
		"enabled": true,
		"type":    "PositiveInt64",
	},
}
