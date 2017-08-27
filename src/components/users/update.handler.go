package users

import (
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/types/apierror"
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
	Name            string `from:"form" json:"name" params:"trim" maxlen:"255"`
	Email           string `from:"form" json:"email" params:"trim,email" maxlen:"255"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim" maxlen:"255"`
	NewPassword     string `from:"form" json:"new_password" params:"trim" maxlen:"255"`

	LastName         *string `from:"form" json:"last_name" params:"trim" maxlen:"50"`
	FirstName        *string `from:"form" json:"first_name" params:"trim" maxlen:"50"`
	PhoneNumber      *string `from:"form" json:"phone_number" params:"trim" maxlen:"255"`
	PublicEmail      *string `from:"form" json:"public_email" params:"trim,email" maxlen:"255"`
	LinkedIn         *string `from:"form" json:"linkedin_custom_url" params:"trim" maxlen:"255"`
	FacebookUsername *string `from:"form" json:"facebook_username" params:"trim" maxlen:"255"`
	TwitterUsername  *string `from:"form" json:"twitter_username" params:"trim" maxlen:"255"`
}

// Update is a HTTP handler used to update a user
func Update(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*UpdateParams)
	currentUser := req.User()

	if !currentUser.IsAdmin && params.ID != currentUser.ID {
		return apierror.NewForbidden()
	}

	profile, err := GetByIDWithProfile(deps.DB, params.ID)
	if err != nil {
		return err
	}

	// To change the email or the password we require the current password
	if params.NewPassword != "" || params.Email != "" {
		if !auth.IsPasswordValid(profile.User.Password, params.CurrentPassword) {
			return apierror.NewUnauthorized()
		}
	}

	// Update the User Object
	if params.Name != "" {
		profile.User.Name = params.Name
	}
	if params.Email != "" {
		profile.User.Email = params.Email
	}
	if params.NewPassword != "" {
		hashedPassword, err := auth.CryptPassword(params.NewPassword)
		if err != nil {
			return err
		}
		profile.User.Password = hashedPassword
	}

	// Update the Profile object
	if params.LastName != nil {
		profile.LastName = params.LastName
		if *profile.LastName == "" {
			profile.LastName = nil
		}
	}
	if params.FirstName != nil {
		profile.FirstName = params.FirstName
		if *profile.FirstName == "" {
			profile.FirstName = nil
		}
	}
	if params.PhoneNumber != nil {
		profile.PhoneNumber = params.PhoneNumber
		if *profile.PhoneNumber == "" {
			profile.PhoneNumber = nil
		}
	}
	if params.PublicEmail != nil {
		profile.PublicEmail = params.PublicEmail
		if *profile.PublicEmail == "" {
			profile.PublicEmail = nil
		}
	}
	if params.LinkedIn != nil {
		profile.LinkedIn = params.LinkedIn
		if *profile.LinkedIn == "" {
			profile.LinkedIn = nil
		}
	}
	if params.FacebookUsername != nil {
		profile.FacebookUsername = params.FacebookUsername
		if *profile.FacebookUsername == "" {
			profile.FacebookUsername = nil
		}
	}
	if params.TwitterUsername != nil {
		profile.TwitterUsername = params.TwitterUsername
		if *profile.TwitterUsername == "" {
			profile.TwitterUsername = nil
		}
	}

	// TODO(melvin): use a transaction
	if err := profile.User.Save(deps.DB); err != nil {
		return err
	}
	if err := profile.Save(deps.DB); err != nil {
		return err
	}

	return req.Response().Ok(profile.ExportPrivate())
}
