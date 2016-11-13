package users

import (
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/router"
)

// HandlerDeleteParams represent the request params accepted by HandlerDelete
type HandlerDeleteParams struct {
	ID              string `from:"url" json:"id" params:"uuid"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
}

// HandlerDelete represent an API handler to remove a user
func HandlerDelete(req *router.Request) {
	params := req.Params.(*HandlerDeleteParams)
	user := req.User

	if params.ID != user.ID {
		req.Error(apierror.NewForbidden())
		return
	}

	if !auth.IsPasswordValid(user.Password, params.CurrentPassword) {
		req.Error(apierror.NewUnauthorized())
		return
	}

	if err := user.Delete(); err != nil {
		req.Error(err)
		return
	}

	req.NoContent()
}