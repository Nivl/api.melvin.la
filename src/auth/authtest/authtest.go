package authtest

import (
	"fmt"
	"testing"

	"github.com/dchest/uniuri"
	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/testhelpers"
)

// NewAuth creates a new user and their session
func NewAuth(t *testing.T) (*auth.User, *auth.Session) {
	user := NewUser(t, nil)
	session := &auth.Session{
		UserID: user.ID,
	}

	if err := session.Create(); err != nil {
		t.Fatal(err)
	}

	testhelpers.SaveModels(t, session)
	return user, session
}

// NewUser creates a new user with "fake" as password
func NewUser(t *testing.T, u *auth.User) *auth.User {
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

	if err := u.Create(); err != nil {
		t.Fatalf("failed to create user: %s", err)
	}

	testhelpers.SaveModels(t, u)
	return u
}
