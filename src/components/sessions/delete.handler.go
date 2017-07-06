package sessions

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db"
)

var deleteEndpoint = &router.Endpoint{
	Verb:    "DELETE",
	Path:    "/sessions/{token}",
	Handler: Delete,
	Guard: &guard.Guard{
		ParamStruct: &DeleteParams{},
		Auth:        guard.LoggedUserAccess,
	},
}

// DeleteParams represent the request params accepted by HandlerDelete
type DeleteParams struct {
	Token           string `from:"url" json:"token" params:"uuid,required"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
}

// Delete represent an API handler to remove a session
func Delete(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*DeleteParams)

	// If a user tries to delete their current session, then no need for the
	// password (that's just a logout)
	if params.Token != req.Session().ID {
		if !auth.IsPasswordValid(req.User().Password, params.CurrentPassword) {
			return httperr.NewUnauthorized()
		}
	}

	var session auth.Session
	stmt := "SELECT * FROM sessions WHERE id=$1 AND deleted_at IS NULL LIMIT 1"
	err := db.Get(deps.DB, &session, stmt, params.Token)
	if err != nil {
		return err
	}

	// We always return a 404 in case of a user error to avoid brute-force
	if session.ID == "" || session.UserID != req.User().ID {
		return httperr.NewNotFound()
	}

	if err := session.Delete(deps.DB); err != nil {
		return err
	}

	req.Response().NoContent()
	return nil
}
