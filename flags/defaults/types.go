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

var permittedFlags = map[string]bool{
	// cloud network vip create
	"cloud_network_vip_create_zone": true,
	// cloud private-parent create
	"cloud_private-parent_create_config-id": true,
	"cloud_private-parent_create_zone":      true,
	// cloud server clone
	"cloud_server_clone_config-id": true,
	// cloud server create
	"cloud_server_create_zone":      true,
	"cloud_server_create_template":  true,
	"cloud_server_create_config-id": true,
	// cloud server resize
	"cloud_server_resize_config-id": true,
	// cloud template restore
	"cloud_template_restore_template": true,

	// network ip-pool create
	"network_ip-pool_create_zone": true,
}
