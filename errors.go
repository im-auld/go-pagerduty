package pagerduty

import "fmt"

// InvalidResourceTypeError is a custom error type.
type InvalidResourceTypeError struct {
	Message string
}

func (e InvalidResourceTypeError) Error() string {
	return e.Message
}

// NewInvalidResourceTypeError creates a new `InvalidResourceTypeError`.
func NewInvalidResourceTypeError(typ APIResourceType) InvalidResourceTypeError {
	msg := fmt.Sprintf("%s is not a known resource type.", typ)
	return InvalidResourceTypeError{Message: msg}
}
