package users

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/security/auth"
)

// GetParams represent the request params accepted by HandlerGet
type GetParams struct {
	ID string `from:"url" json:"id" params:"uuid"`
}

// Get represent an API handler to get a user
func Get(req *router.Request, deps *router.Dependencies) error {
	params := req.Params.(*GetParams)

	user, err := auth.GetUser(deps.DB, params.ID)
	if err != nil {
		return err
	}
	if user == nil {
		return httperr.NewNotFound()
	}

	// if a user asks for their own data, we return as much as possible
	if req.User != nil && req.User.ID == user.ID {
		req.Ok(NewPrivatePayload(user))
		return nil
	}
	req.Ok(NewPayload(user))
	return nil
}
