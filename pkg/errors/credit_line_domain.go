package errors

import (
	"errors"
)

var (
	// ErrInvalidFoundingType is returned when the foundingType is invalid
	ErrInvalidFoundingType = errors.New("invalid foundingType")
)
