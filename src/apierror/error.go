package apierror

import (
	"errors"
	"fmt"
	"net/http"
)

// Error represents an error with a code attached.
type Error interface {
	error
	Code() int
}

type ApiError struct {
	error
	ErrorCode int
}

func (err *ApiError) Code() int {
	if err == nil {
		return http.StatusInternalServerError
	}

	return err.ErrorCode
}

// NewError returns an error with an associated code
func NewError(code int, message string, args ...interface{}) error {
	fullMessage := fmt.Sprintf(message, args...)
	return &ApiError{errors.New(fullMessage), code}
}

// NewServerError returns an Internal Error.
func NewServerError(message string, args ...interface{}) error {
	return NewError(http.StatusInternalServerError, message, args)
}

// NewBadRequest returns an error caused by a user. Example: A missing param
func NewBadRequest(message string, args ...interface{}) error {
	return NewError(http.StatusBadRequest, message, args)
}

// NewConflict returns an error caused by a conflict with the current state
// of the app. Example: A duplicate slug
func NewConflict(message string, args ...interface{}) error {
	return NewError(http.StatusConflict, message, args)
}

// NewUnauthorized returns an error caused by a anonymous user trying to access
// a protected resource
func NewUnauthorized() error {
	return NewError(http.StatusUnauthorized, "")
}

// NewForbidden returns an error caused by a user trying to access
// a protected resource
func NewForbidden() error {
	return NewError(http.StatusForbidden, "")
}

// NewNotFound returns an error caused by a user trying to access
// a resource that does not exists
func NewNotFound() error {
	return NewError(http.StatusNotFound, "")
}
