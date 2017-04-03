package users

import (
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/security/auth"
)

// AddParams represents the params accepted by Add to create a new user
type AddParams struct {
	Name     string `from:"form" json:"name" params:"required,trim"`
	Email    string `from:"form" json:"email" params:"required,trim"`
	Password string `from:"form" json:"password" params:"required,trim"`
}

// Add is a HTTP handler used to add a new user
func Add(req *router.Request) error {
	params := req.Params.(*AddParams)

	encryptedPassword, err := auth.CryptPassword(params.Password)
	if err != nil {
		return err
	}

	user := &auth.User{
		Name:     params.Name,
		Email:    params.Email,
		Password: encryptedPassword,
	}

	if err := user.Save(); err != nil {
		return err
	}

	req.Created(NewPrivatePayload(user))
	return nil
}
