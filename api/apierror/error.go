package apierror

import (
	"errors"
	"fmt"
	"net/http"
)

// Error represents an error with a code attached.
type Error struct {
	error
	Code int
}

// NewError returns an error with an associated code
func NewError(code int, message string, args ...interface{}) error {
	fullMessage := fmt.Sprintf(message, args...)

	return Error{errors.New(fullMessage), code}
}

// NewServerError returns an Internal Error.
func NewServerError(message string, args ...interface{}) error {
	return NewError(http.StatusInternalServerError, message, args)
}

// NewBadRequest returns an error caused by a user. Example: A missing param
func NewBadRequest(message string, args ...interface{}) error {
	return NewError(http.StatusInternalServerError, message, args)
}

// NewConflict returns an error caused by a conflict with the current state
// of the app. Example: A duplicate slug
func NewConflict(message string, args ...interface{}) error {
	return NewError(http.StatusConflict, message, args)
}
