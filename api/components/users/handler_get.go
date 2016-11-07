package users

import (
	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/router"
)

// HandlerGetParams represent the request params accepted by HandlerGet
type HandlerGetParams struct {
	ID string `from:"url" json:"id" params:"uuid"`
}

// HandlerGet represent an API handler to get a user
func HandlerGet(req *router.Request) {
	params := req.Params.(*HandlerGetParams)

	user, err := auth.GetUser(params.ID)
	if err != nil {
		req.Error(err)
		return
	}
	if user == nil {
		req.Error(apierror.NewNotFound())
		return
	}

	// if a user asks for their own data, we return as much as possible
	if req.User != nil && req.User.ID == user.ID {
		req.Ok(NewPrivatePayload(user))
		return
	}
	req.Ok(NewPublicPayload(user))
}
