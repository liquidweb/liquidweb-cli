package errorTypes

import (
	"errors"
)

var LwCliInputError = errors.New("Invalid input; missing required paramater")
var LwApiUnexpectedResponseStructure = errors.New("Unexpected API response structure when calling method")
var UnknownTerminal = errors.New("unknown terminal")
var MergeConfigError = errors.New("error merging configuration")
