package users

import (
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/router"
)

// HandlerGetParams represent the request params accepted by HandlerGet
type HandlerGetParams struct {
	ID string `from:"url" json:"id" params:"uuid"`
}

// HandlerGet represent an API handler to get a user
func HandlerGet(req *router.Request) error {
	params := req.Params.(*HandlerGetParams)

	user, err := auth.GetUser(params.ID)
	if err != nil {
		return err
	}
	if user == nil {
		return apierror.NewNotFound()
	}

	// if a user asks for their own data, we return as much as possible
	if req.User != nil && req.User.ID == user.ID {
		req.Ok(NewPrivatePayload(user))
		return nil
	}
	req.Ok(NewPublicPayload(user))
	return nil
}
