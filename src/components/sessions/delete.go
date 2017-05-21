package sessions

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db"
)

// DeleteParams represent the request params accepted by HandlerDelete
type DeleteParams struct {
	Token           string `from:"url" json:"token" params:"uuid"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
}

// Delete represent an API handler to remove a session
func Delete(req *router.Request) error {
	params := req.Params.(*DeleteParams)

	// If a user tries to delete their current session, then no need for the
	// password (that's just a logout)
	if params.Token != req.SessionUsed.ID {
		if !auth.IsPasswordValid(req.User.Password, params.CurrentPassword) {
			return httperr.NewUnauthorized()
		}
	}

	var session auth.Session
	stmt := "SELECT * FROM sessions WHERE id=$1 AND deleted_at IS NULL LIMIT 1"
	err := db.Get(&session, stmt, params.Token)
	if err != nil {
		return err
	}

	// We always return a 404 in case of a user error to avoid brute-force
	if session.ID == "" || session.UserID != req.User.ID {
		return httperr.NewNotFound()
	}

	if err := session.Delete(); err != nil {
		return err
	}

	req.NoContent()
	return nil
}
