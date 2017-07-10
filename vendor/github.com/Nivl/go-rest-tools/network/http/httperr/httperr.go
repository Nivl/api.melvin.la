package httperr

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/lib/pq"
)

// Error represents an error with a code attached.
type Error interface {
	error

	// Code return the HTTP code of the error
	Code() int

	// Field returns the http/sql/etc. field associated to the error
	Field() string

	// Origin returns the original error if there's one
	Origin() error
}

// Convert takes an error an turns it into an HTTPError
func Convert(e error) *HTTPError {
	err, casted := e.(*HTTPError)
	if !casted {
		err = NewServerError(e.Error())
		err.ErrorOrigin = e
	}

	return err
}

// HTTPError represents an error with an HTTP code
type HTTPError struct {
	error
	ErrorCode   int
	ErrorField  string
	ErrorOrigin error
}

// Code returns the HTTP code associated to the error
func (err *HTTPError) Code() int {
	if err == nil {
		return http.StatusInternalServerError
	}

	return err.ErrorCode
}

// Field returns the HTTP param associated to the error
func (err *HTTPError) Field() string {
	return err.ErrorField
}

// Origin returns the original error
func (err *HTTPError) Origin() error {
	return err.ErrorOrigin
}

// NewFromSQL returns an http error based on a pq.Error
// the provided error will be returned if it's not a pq.Error instance,
// or if the error cannot be matched to
func NewFromSQL(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case db.ErrDup:
			// because it's a constraint issue, the column name won't be stored in
			// pqErr.Column. Fortunately we can find it in detail.
			// Example of detail: "Key (name)=(Google) already exists."
			r := regexp.MustCompile(`^Key \(([a-z_]+)\).*\.$`)
			matches := r.FindStringSubmatch(pqErr.Detail)
			fieldName := "unknown"
			if len(matches) > 1 {
				fieldName = matches[1]
			}
			return NewConflict(fieldName)
		default:
			return err
		}
	}
	return err
}

// NewError returns an error with an associated code
func NewError(code int, field string, message string, args ...interface{}) *HTTPError {
	fullMessage := fmt.Sprintf(message, args...)
	return &HTTPError{errors.New(fullMessage), code, field, nil}
}

// NewServerError returns an Internal Error.
func NewServerError(message string, args ...interface{}) *HTTPError {
	return NewError(http.StatusInternalServerError, "", message, args...)
}

// NewBadRequest returns an error caused by a user. Example: A missing param
func NewBadRequest(field string, message string, args ...interface{}) *HTTPError {
	return NewError(http.StatusBadRequest, field, message, args...)
}

// NewConflict returns an error caused by a conflict with the current state
// of the app. Example: A duplicate slug
func NewConflict(field string) *HTTPError {
	return NewError(http.StatusConflict, field, "already exists")
}

// NewConflictR returns an error caused by a conflict with the current state
// of the app. A reason is sent back to the user.
func NewConflictR(field string, message string, args ...interface{}) *HTTPError {
	return NewError(http.StatusConflict, field, message, args...)
}

// NewUnauthorized returns an error caused by a anonymous user trying to access
// a protected resource
func NewUnauthorized() *HTTPError {
	return NewUnauthorizedR(http.StatusText(http.StatusUnauthorized))
}

// NewUnauthorizedR returns an error caused by a anonymous user trying to access
// a protected resource. A reason is sent back to the user.
func NewUnauthorizedR(reason string) *HTTPError {
	return NewError(http.StatusUnauthorized, "", reason)
}

// NewForbidden returns an error caused by a user trying to access
// a protected resource.
func NewForbidden() *HTTPError {
	return NewForbiddenR(http.StatusText(http.StatusForbidden))
}

// NewForbiddenR returns an error caused by a user trying to access
// a protected resource. A reason is sent back to the user.
func NewForbiddenR(reason string) *HTTPError {
	return NewError(http.StatusForbidden, "", reason)
}

// NewNotFound returns an error caused by a user trying to access
// a resource that does not exists
func NewNotFound() *HTTPError {
	return NewNotFoundR(http.StatusText(http.StatusNotFound))
}

// NewNotFoundR returns an error caused by a user trying to access
// a resource that does not exists. A reason is sent back to the user.
func NewNotFoundR(reason string) *HTTPError {
	return NewError(http.StatusNotFound, "", reason)
}

// NewNotFoundField returns an error caused by a user trying to access
// a resource that does not exists. A reason is sent back to the user.
func NewNotFoundField(field string, reason string) *HTTPError {
	return NewError(http.StatusNotFound, field, reason)
}
