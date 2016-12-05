package users

import "github.com/melvin-laplanche/ml-api/src/router"
import "github.com/melvin-laplanche/ml-api/src/auth"

type HandlerAddParams struct {
	Name     string `from:"form" json:"name" params:"required,trim"`
	Email    string `from:"form" json:"email" params:"required,trim"`
	Password string `from:"form" json:"password" params:"required,trim"`
}

func HandlerAdd(req *router.Request) error {
	params := req.Params.(*HandlerAddParams)

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
