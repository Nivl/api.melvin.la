package auth

import (
	"errors"
	"time"

	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/Nivl/api.melvin.la/api/app"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func QuerySessions() *mgo.Collection {
	return app.GetContext().DB.C("sessions")
}

// Session is a structure representing a session that can be saved in the database
type Session struct {
	ID        bson.ObjectId `bson:"_id"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
	IsDeleted bool          `bson:"is_deleted"`

	UserID bson.ObjectId `bson:"user_id"`
}

// NewSessionFromStrings returns a Session using the provided strings
func NewSessionFromStrings(sessionID string, userID string) *Session {
	sess := &Session{}

	if bson.IsObjectIdHex(sessionID) {
		sess.ID = bson.ObjectIdHex(sessionID)
	}

	if bson.IsObjectIdHex(userID) {
		sess.UserID = bson.ObjectIdHex(userID)
	}

	return sess
}

// Exists check if a session exists in the database
func (s *Session) Exists() (bool, error) {
	if s == nil {
		return false, apierror.NewServerError("session is nil")
	}

	if s.UserID == "" {
		return false, apierror.NewServerError("user id required")
	}

	// Deleted sessions should be explicitly checked
	if s.IsDeleted {
		return false, nil
	}

	count, err := QuerySessions().Find(s).Count()
	if err != nil {
		return false, err
	}

	return (count > 0), nil
}

// Save is an alias for Create
func (s *Session) Save() error {
	if s == nil {
		return apierror.NewServerError("session is nil")
	}

	return s.Create()
}

// Create persists a session in the database
func (s *Session) Create() error {
	if s == nil {
		return apierror.NewServerError("session is nil")
	}

	if s.ID != "" {
		return apierror.NewServerError("sessions cannot be updated")
	}

	if s.UserID == "" {
		return apierror.NewServerError("cannot create a session without a user id")
	}

	s.ID = bson.NewObjectId()
	s.CreatedAt = time.Now()
	return QuerySessions().Insert(s)
}

func (s *Session) FullyDelete() error {
	if s == nil {
		return errors.New("session not instanced")
	}

	if s.ID == "" {
		return errors.New("session has not been saved")
	}

	return QuerySessions().RemoveId(s.ID)
}
