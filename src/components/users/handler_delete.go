package users

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/security/auth"
)

// HandlerDeleteParams represent the request params accepted by HandlerDelete
type HandlerDeleteParams struct {
	ID              string `from:"url" json:"id" params:"uuid"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
}

// HandlerDelete represent an API handler to remove a user
func HandlerDelete(req *router.Request) error {
	params := req.Params.(*HandlerDeleteParams)
	user := req.User

	if params.ID != user.ID {
		return httperr.NewForbidden()
	}

	if !auth.IsPasswordValid(user.Password, params.CurrentPassword) {
		return httperr.NewUnauthorized()
	}

	if err := user.Delete(); err != nil {
		return err
	}

	req.NoContent()
	return nil
}
