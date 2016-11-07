package sessions

import (
	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/db"
	"github.com/Nivl/api.melvin.la/api/router"
)

// HandlerAddParams represent the request params accepted by HandlerAdd
type HandlerAddParams struct {
	Email    string `from:"form" json:"email" params:"required,trim"`
	Password string `from:"form" json:"password" params:"required,trim"`
}

// HandlerAdd represents an API handler to create a new user session
func HandlerAdd(req *router.Request) {
	params := req.Params.(*HandlerAddParams)

	var user auth.User
	stmt := "SELECT * FROM users WHERE email=$1 LIMIT 1"
	err := db.Get(&user, stmt, params.Email)
	if err != nil {
		req.Error(err)
		return
	}

	if user.ID == "" || !auth.IsPasswordValid(user.Password, params.Password) {
		req.Error(apierror.NewBadRequest("Bad email/password"))
		return
	}

	s := &auth.Session{
		UserID: user.ID,
	}
	if err := s.Save(); err != nil {
		req.Error(err)
		return
	}

	req.Created(NewPayloadFromModel(s))
}
