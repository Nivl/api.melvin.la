package users

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/security/auth"
)

// DeleteParams represent the request params accepted by HandlerDelete
type DeleteParams struct {
	ID              string `from:"url" json:"id" params:"uuid"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
}

// Delete represent an API handler to remove a user
func Delete(req *router.Request, deps *router.Dependencies) error {
	params := req.Params.(*DeleteParams)
	user := req.User

	if params.ID != user.ID {
		return httperr.NewForbidden()
	}

	if !auth.IsPasswordValid(user.Password, params.CurrentPassword) {
		return httperr.NewUnauthorized()
	}

	if err := user.Delete(deps.DB); err != nil {
		return err
	}

	req.NoContent()
	return nil
}
