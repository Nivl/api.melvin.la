package users

import (
	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/router"
	"gopkg.in/mgo.v2/bson"
)

// HandlerGetParams represent the request params accepted by HandlerGet
type HandlerGetParams struct {
	ID string `from:"url" json:"id"`
}

// HandlerGet represent an API handler to get a user
func HandlerGet(req *router.Request) {
	params := req.Params.(*HandlerGetParams)

	// ID is not valid
	if !bson.IsObjectIdHex(params.ID) {
		req.Error(apierror.NewNotFound())
		return
	}

	user, err := auth.GetUser(bson.ObjectIdHex(params.ID))
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
