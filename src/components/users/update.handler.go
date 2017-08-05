package users

import (
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/security/auth"
)

var updateEndpoint = &router.Endpoint{
	Verb:    "PATCH",
	Path:    "/users/{id}",
	Handler: Update,
	Guard: &guard.Guard{
		ParamStruct: &UpdateParams{},
		Auth:        guard.LoggedUserAccess,
	},
}

// UpdateParams represents the params accepted Update to update a user
type UpdateParams struct {
	ID              string `from:"url" json:"id"  params:"uuid,required"`
	Name            string `from:"form" json:"name" params:"trim"`
	Email           string `from:"form" json:"email" params:"trim"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
	NewPassword     string `from:"form" json:"new_password" params:"trim"`
}

// Update is a HTTP handler used to update a user
func Update(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*UpdateParams)
	user := req.User()

	if params.ID != user.ID {
		return apierror.NewForbidden()
	}

	// To change the email or the password we require the current password
	if params.NewPassword != "" || params.Email != "" {
		if !auth.IsPasswordValid(user.Password, params.CurrentPassword) {
			return apierror.NewUnauthorized()
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

	return req.Response().Ok(NewPrivatePayload(user))
}
