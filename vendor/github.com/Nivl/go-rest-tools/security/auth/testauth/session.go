package testauth

import (
	"testing"

	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
)

// NewAuth creates a new user and their session
func NewAuth(t *testing.T, q db.DB) (*auth.User, *auth.Session) {
	user := NewPersistedUser(t, q, nil)
	session := NewPersistedSession(t, q, user)
	return user, session
}

// NewAdminAuth creates a new admin and their session
func NewAdminAuth(t *testing.T, q db.DB) (*auth.User, *auth.Session) {
	user := NewPersistedUser(t, q, &auth.User{IsAdmin: true})
	session := NewPersistedSession(t, q, user)
	return user, session
}

// NewSession creates a non-persisted session for the given user
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

// NewPersistedSession creates and persists a new session for the given user
func NewPersistedSession(t *testing.T, q db.DB, user *auth.User) *auth.Session {
	session := &auth.Session{
		UserID: user.ID,
	}

	if err := session.Create(q); err != nil {
		t.Fatal(err)
	}

	lifecycle.SaveModels(t, session)
	return session
}
