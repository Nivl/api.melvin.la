package users

import (
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/router"
)

type HandlerUpdateParams struct {
	ID              string `from:"url" json:"id"  params:"uuid"`
	Name            string `from:"form" json:"name" params:"trim"`
	Email           string `from:"form" json:"email" params:"trim"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
	NewPassword     string `from:"form" json:"new_password" params:"trim"`
}

func HandlerUpdate(req *router.Request) {
	params := req.Params.(*HandlerUpdateParams)
	user := req.User

	if params.ID != user.ID {
		req.Error(apierror.NewForbidden())
		return
	}

	// To change the email or the password we require the current password
	if params.NewPassword != "" || params.Email != "" {
		if !auth.IsPasswordValid(user.Password, params.CurrentPassword) {
			req.Error(apierror.NewUnauthorized())
			return
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
			req.Error(err)
			return
		}
		user.Password = hashedPassword
	}

	if err := user.Save(); err != nil {
		req.Error(err)
		return
	}

	req.Ok(NewPrivatePayload(user))
}
