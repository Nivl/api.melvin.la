package guard

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/security/auth"
)

// RouteAuth represents a middleware used to allow/block the access to an endpoint
type RouteAuth func(*auth.User) error

// LoggedUserAccess is a auth middleware that filters out anonymous users
func LoggedUserAccess(u *auth.User) error {
	if u == nil || u.ID == "" {
		return httperr.NewUnauthorized()
	}
	return nil
}

// AdminAccess is a auth middleware that filters out non admin users
func AdminAccess(u *auth.User) error {
	if err := LoggedUserAccess(u); err != nil {
		return httperr.NewUnauthorized()
	}
	if !u.IsAdmin {
		return httperr.NewForbidden()
	}
	return nil
}
