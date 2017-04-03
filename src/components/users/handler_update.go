package users

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/security/auth"
)

type HandlerUpdateParams struct {
	ID              string `from:"url" json:"id"  params:"uuid"`
	Name            string `from:"form" json:"name" params:"trim"`
	Email           string `from:"form" json:"email" params:"trim"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
	NewPassword     string `from:"form" json:"new_password" params:"trim"`
}

func HandlerUpdate(req *router.Request) error {
	params := req.Params.(*HandlerUpdateParams)
	user := req.User

	if params.ID != user.ID {
		return httperr.NewForbidden()
	}

	// To change the email or the password we require the current password
	if params.NewPassword != "" || params.Email != "" {
		if !auth.IsPasswordValid(user.Password, params.CurrentPassword) {
			return httperr.NewUnauthorized()
		}
	}

	if params.Name != "" {
		user.Name = params.Name
	}

	if params.Email != "" {
		user.Email = params.Email
	}

	if params.NewPassword != "" {
		hashedPassword, err := auth.CryptPassword(params.NewPassword)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	if err := user.Save(); err != nil {
		return err
	}

	req.Ok(NewPrivatePayload(user))
	return nil
}
