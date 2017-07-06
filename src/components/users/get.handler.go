package users

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/security/auth"
)

var getEndpoint = &router.Endpoint{
	Verb:    "GET",
	Path:    "/users/{id}",
	Handler: Get,
	Guard: &guard.Guard{
		ParamStruct: &GetParams{},
	},
}

// GetParams represent the request params accepted by HandlerGet
type GetParams struct {
	ID string `from:"url" json:"id" params:"uuid,required"`
}

// Get represent an API handler to get a user
func Get(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*GetParams)

	// if a user asks for their own data, we don't need to query the DB
	if req.User() != nil && req.User().ID == params.ID {
		return req.Response().Ok(NewPrivatePayload(req.User()))
	}

	user, err := auth.GetUser(deps.DB, params.ID)
	if err != nil {
		return err
	}
	if user == nil {
		return httperr.NewNotFound()
	}

	return req.Response().Ok(NewPayload(user))
}
