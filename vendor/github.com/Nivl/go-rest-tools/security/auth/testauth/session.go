package testauth

import (
	"testing"

	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db"
)

// NewAuth creates a new user and their session
func NewAuth(t *testing.T, q db.DB) (*auth.User, *auth.Session) {
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
func NewAdminAuth(t *testing.T, q db.DB) (*auth.User, *auth.Session) {
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

// NewSession creates a new session for the given user
func NewSession(t *testing.T, q db.DB, user *auth.User) *auth.Session {
	session := &auth.Session{
		UserID: user.ID,
	}

	if err := session.Create(q); err != nil {
		t.Fatal(err)
	}

	lifecycle.SaveModels(t, session)
	return session
}
