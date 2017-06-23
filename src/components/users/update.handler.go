package users

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/security/auth"
)

// UpdateParams represents the params accepted Update to update a user
type UpdateParams struct {
	ID              string `from:"url" json:"id"  params:"uuid"`
	Name            string `from:"form" json:"name" params:"trim"`
	Email           string `from:"form" json:"email" params:"trim"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
	NewPassword     string `from:"form" json:"new_password" params:"trim"`
}

// Update is a HTTP handler used to update a user
func Update(req *router.Request, deps *router.Dependencies) error {
	params := req.Params.(*UpdateParams)
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

	if err := user.Save(deps.DB); err != nil {
		return err
	}

	req.Ok(NewPrivatePayload(user))
	return nil
}
