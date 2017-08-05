package testauth

import (
	"fmt"
	"testing"

	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/dchest/uniuri"
)

// NewUser creates a new user with "fake" as password
func NewUser(t *testing.T, q db.DB, u *auth.User) *auth.User {
	if u == nil {
		u = &auth.User{}
	}

	if u.Email == "" {
		u.Email = fmt.Sprintf("fake+%s@melvin.la", uniuri.New())
	}

	if u.Name == "" {
		u.Name = "Fake Account"
	}

	if u.Password == "" {
		var err error
		u.Password, err = auth.CryptPassword("fake")
		if err != nil {
			t.Fatalf("failed to create password: %s", err)
		}
	}

	if err := u.Create(q); err != nil {
		t.Fatalf("failed to create user: %s", err)
	}

	lifecycle.SaveModels(t, u)
	return u
}
