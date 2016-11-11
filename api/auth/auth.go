package auth

import (
	"testing"

	"github.com/Nivl/api.melvin.la/api/app"
	"github.com/jmoiron/sqlx"
)

func sql() *sqlx.DB {
	return app.GetContext().SQL
}

// NewTestAuth creates a new user and their session
func NewTestAuth(t *testing.T) (*User, *Session) {
	user := NewTestUser(t, nil)
	session := &Session{
		UserID: user.ID,
	}

	if err := session.Create(); err != nil {
		t.Fatal(err)
	}

	return user, session
}
