package sessions

import (
	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/db"
	"github.com/Nivl/api.melvin.la/api/router"
)

// HandlerDeleteParams represent the request params accepted by HandlerDelete
type HandlerDeleteParams struct {
	Token           string `from:"url" json:"token"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
}

// HandlerDelete represent an API handler to remove a session
func HandlerDelete(req *router.Request) {
	params := req.Params.(*HandlerDeleteParams)

	if !auth.IsPasswordValid(req.User.Password, params.CurrentPassword) {
		req.Error(apierror.NewUnauthorized())
		return
	}

	var session auth.Session
	stmt := "SELECT * FROM sessions WHERE uuid=$1 AND deleted_at IS NULL LIMIT 1"
	err := db.Get(&session, stmt, params.Token)
	if err != nil {
		req.Error(err)
		return
	}

	// We always return a 404 in case of a user error to avoid brute-force
	if session.UUID == "" || session.UserUUID != req.User.UUID {
		req.Error(apierror.NewNotFound())
		return
	}

	if err := session.Delete(); err != nil {
		req.Error(err)
		return
	}

	req.NoContent()
}
