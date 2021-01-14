package defaults

import (
	"errors"
)

var ErrorForbiddenFlag = errors.New("is a forbidden default flag")
var ErrorInvalidFlagName = errors.New("the given flag name is invalid")
var ErrorFileKeyMissing = errors.New("flag defaults file key is missing")
var ErrorUnwritable = errors.New("flag defaults cannot be written")
var ErrorUnreadable = errors.New("flag defaults cannot be read")
var ErrorNotFound = errors.New("flag default not found")
