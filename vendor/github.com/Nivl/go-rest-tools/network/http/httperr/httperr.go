package httperr

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

// Convert takes an error an turns it into an HTTPError
func Convert(e error) *HTTPError {
	err, casted := e.(*HTTPError)
	if !casted {
		err = NewServerError(e.Error()).(*HTTPError)
	}

	return err
}

// HTTPError represents an error with an HTTP code
type HTTPError struct {
	error
	ErrorCode int
}

// Code returns the HTTP code associated to the error
func (err *HTTPError) Code() int {
	if err == nil {
		return http.StatusInternalServerError
	}

	return err.ErrorCode
}

// NewError returns an error with an associated code
func NewError(code int, message string, args ...interface{}) Error {
	fullMessage := fmt.Sprintf(message, args...)
	return &HTTPError{errors.New(fullMessage), code}
}

// NewServerError returns an Internal Error.
func NewServerError(message string, args ...interface{}) Error {
	return NewError(http.StatusInternalServerError, message, args...)
}

// NewBadRequest returns an error caused by a user. Example: A missing param
func NewBadRequest(message string, args ...interface{}) Error {
	return NewError(http.StatusBadRequest, message, args...)
}

// NewConflict returns an error caused by a conflict with the current state
// of the app. Example: A duplicate slug
func NewConflict(message string, args ...interface{}) Error {
	return NewError(http.StatusConflict, message, args...)
}

// NewUnauthorized returns an error caused by a anonymous user trying to access
// a protected resource
func NewUnauthorized() Error {
	return NewUnauthorizedR(http.StatusText(http.StatusUnauthorized))
}

// NewUnauthorizedR returns an error caused by a anonymous user trying to access
// a protected resource. A reason is sent back to the user.
func NewUnauthorizedR(reason string) Error {
	return NewError(http.StatusUnauthorized, reason)
}

// NewForbidden returns an error caused by a user trying to access
// a protected resource.
func NewForbidden() Error {
	return NewForbiddenR(http.StatusText(http.StatusForbidden))
}

// NewForbiddenR returns an error caused by a user trying to access
// a protected resource. A reason is sent back to the user.
func NewForbiddenR(reason string) Error {
	return NewError(http.StatusForbidden, reason)
}

// NewNotFound returns an error caused by a user trying to access
// a resource that does not exists
func NewNotFound() Error {
	return NewNotFoundR(http.StatusText(http.StatusNotFound))
}

// NewNotFoundR returns an error caused by a user trying to access
// a resource that does not exists. A reason is sent back to the user.
func NewNotFoundR(reason string) Error {
	return NewError(http.StatusNotFound, reason)
}
