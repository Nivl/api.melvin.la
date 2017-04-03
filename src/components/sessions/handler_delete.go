package sessions

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db"
)

// HandlerDeleteParams represent the request params accepted by HandlerDelete
type HandlerDeleteParams struct {
	Token           string `from:"url" json:"token" params:"uuid"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
}

// HandlerDelete represent an API handler to remove a session
func HandlerDelete(req *router.Request) error {
	params := req.Params.(*HandlerDeleteParams)

	if !auth.IsPasswordValid(req.User.Password, params.CurrentPassword) {
		return httperr.NewUnauthorized()
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
