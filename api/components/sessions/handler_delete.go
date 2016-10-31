package sessions

import (
	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/router"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// HandlerDeleteParams represent the request params accepted by HandlerDelete
type HandlerDeleteParams struct {
	ID              string `from:"url" json:"id"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
}

// HandlerDelete represent an API handler to remove a session
func HandlerDelete(req *router.Request) {
	params := req.Params.(*HandlerDeleteParams)

	if !auth.IsPasswordValid(req.User.Password, params.CurrentPassword) {
		req.Error(apierror.NewUnauthorized())
		return
	}

	if !bson.IsObjectIdHex(params.ID) {
		req.Error(apierror.NewNotFound())
		return
	}

	toFind := bson.M{
		"_id":        bson.ObjectIdHex(params.ID),
		"is_deleted": false,
	}

	var session auth.Session
	err := auth.QuerySessions().Find(&toFind).One(&session)
	if err != nil && err != mgo.ErrNotFound {
		req.Error(err)
		return
	}

	// We always return a 404 in case of a user error to avoid
	if session.ID == "" || session.UserID != req.User.ID {
		req.Error(apierror.NewNotFound())
		return
	}

	if err := session.Delete(); err != nil {
		req.Error(err)
		return
	}

	req.NoContent()
}
