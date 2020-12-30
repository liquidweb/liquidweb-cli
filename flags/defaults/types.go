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
			slice = append(slice, fmt.Sprintf("\tFlag: %s\n", flag))
			slice = append(slice, fmt.Sprintf("\t\tValue: %+v\n", value))
		}
	}

	return strings.Join(slice[:], "")
}

var permittedFlags = map[string]interface{}{
	"zone":     true,
	"template": true,
}
