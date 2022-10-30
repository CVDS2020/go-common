package errors

import (
	"errors"
)

var (
	New    = errors.New
	Unwrap = errors.Unwrap
	As     = errors.As
	Is     = errors.Is
)
