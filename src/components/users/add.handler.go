package users

import (
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/security/auth"
)

var addEndpoint = &router.Endpoint{
	Verb:    "POST",
	Path:    "/users",
	Handler: Add,
	Guard: &guard.Guard{
		ParamStruct: &AddParams{},
	},
}

// AddParams represents the params accepted by Add to create a new user
type AddParams struct {
	Name     string `from:"form" json:"name" params:"required,trim"`
	Email    string `from:"form" json:"email" params:"required,trim"`
	Password string `from:"form" json:"password" params:"required,trim"`
}

// Add is a HTTP handler used to add a new user
func Add(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*AddParams)

	encryptedPassword, err := auth.CryptPassword(params.Password)
	if err != nil {
		return err
	}

	user := &auth.User{
		Name:     params.Name,
		Email:    params.Email,
		Password: encryptedPassword,
	}

	if err := user.Create(deps.DB); err != nil {
		return err
	}

	profile := &Profile{User: user}
	return req.Response().Created(profile.ExportPrivate())
}
