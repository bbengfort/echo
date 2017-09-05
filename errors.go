package echo

import (
	"errors"
	"fmt"
)

//===========================================================================
// Some standard errors that may be thrown.
//===========================================================================

// Standard errors for primary operations.
var (
	ErrNotImplemented = errors.New("functionality not implemented yet")
)

//===========================================================================
// Error wraps other library errors for ease of logging.
//===========================================================================

// Error wraps other library errors and provides easier logging.
type Error struct {
	msg string // The alia error message
	err error  // The wrapped error if any
}

// WrapError creates a new wrapped error message with the format string.
func WrapError(format string, err error, a ...interface{}) *Error {
	return &Error{
		msg: fmt.Sprintf(format, a...),
		err: err,
	}
}

// Error prefixes the message to the internal error string
func (e *Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	}

	return e.msg
}

// String returns the error message
func (e *Error) String() string {
	return e.Error()
}
