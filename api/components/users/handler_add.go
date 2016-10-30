package users

import "github.com/Nivl/api.melvin.la/api/router"
import "github.com/Nivl/api.melvin.la/api/auth"

type HandlerAddParams struct {
	Name     string `from:"form" json:"name" params:"required,trim"`
	Email    string `from:"form" json:"email" params:"required,trim"`
	Password string `from:"form" json:"password" params:"required,trim"`
}

func HandlerAdd(req *router.Request) {
	params := req.Params.(*HandlerAddParams)

	encryptedPassword, err := auth.CryptPassword(params.Password)
	if err != nil {
		req.Error(err)
		return
	}

	user := &auth.User{
		Name:     params.Name,
		Email:    params.Email,
		Password: encryptedPassword,
	}

	if err := user.Save(); err != nil {
		req.Error(err)
		return
	}

	req.Created(NewPrivatePayload(user))
}
