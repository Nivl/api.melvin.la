package testdata

import (
	"testing"

	"github.com/Nivl/go-rest-tools/primitives/models/lifecycle"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/jmoiron/sqlx"
)

// NewAuth creates a new user and their session
func NewAuth(t *testing.T, q *sqlx.DB) (*auth.User, *auth.Session) {
	user := NewUser(t, q, nil)
	session := &auth.Session{
		UserID: user.ID,
	}

	if err := session.Create(q); err != nil {
		t.Fatal(err)
	}

	lifecycle.SaveModels(t, session)
	return user, session
}

// NewAdminAuth creates a new admin and their session
func NewAdminAuth(t *testing.T, q *sqlx.DB) (*auth.User, *auth.Session) {
	user := NewUser(t, q, &auth.User{IsAdmin: true})
	session := &auth.Session{
		UserID: user.ID,
	}

	if err := session.Create(q); err != nil {
		t.Fatal(err)
	}

	lifecycle.SaveModels(t, session)
	return user, session
}
