package testusers

import (
	"testing"

	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/Nivl/go-rest-tools/types/ptrs"
	"github.com/dchest/uniuri"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	uuid "github.com/satori/go.uuid"
)

// NewProfile returns a non persisted profile and user
func NewProfile() *users.Profile {
	u := testauth.NewUser()
	return &users.Profile{
		ID:               uuid.NewV4().String(),
		LastName:         ptrs.NewString(uniuri.New()),
		FirstName:        ptrs.NewString(uniuri.New()),
		LinkedIn:         ptrs.NewString(uniuri.New()),
		FacebookUsername: ptrs.NewString(uniuri.New()),
		TwitterUsername:  ptrs.NewString(uniuri.New()),
		PublicEmail:      ptrs.NewString(uniuri.New() + "@domain.tld"),
		PhoneNumber:      ptrs.NewString("+1 (123) 456 7890"),
		Picture:          ptrs.NewString("http://placehold.it/60x60"),
		UserID:           u.ID,
		User:             u,
	}
}

// NewPersistedProfile creates and persists a new user and user profile with "fake" as password
func NewPersistedProfile(t *testing.T, q db.DB, p *users.Profile) *users.Profile {
	if p == nil {
		p = &users.Profile{}
	}

	p.User = testauth.NewPersistedUser(t, q, p.User)
	p.UserID = p.User.ID

	if p.LastName == nil {
		p.LastName = ptrs.NewString(uniuri.New())
	}
	if p.FirstName == nil {
		p.FirstName = ptrs.NewString(uniuri.New())
	}
	if p.LinkedIn == nil {
		p.LinkedIn = ptrs.NewString(uniuri.New())
	}
	if p.FacebookUsername == nil {
		p.FacebookUsername = ptrs.NewString(uniuri.New())
	}
	if p.TwitterUsername == nil {
		p.TwitterUsername = ptrs.NewString(uniuri.New())
	}
	if p.PublicEmail == nil {
		p.PublicEmail = ptrs.NewString(uniuri.New() + "@domain.tld")
	}
	if p.Picture == nil {
		p.Picture = ptrs.NewString("http://placehold.it/60x60")
	}
	if p.PhoneNumber == nil {
		p.PhoneNumber = ptrs.NewString("+1 (123) 456 7890")
	}

	if err := p.Create(q); err != nil {
		t.Fatalf("failed to create user: %s", err)
	}

	lifecycle.SaveModels(t, p)
	return p
}

// NewAuth creates a new user and their session
func NewAuth(t *testing.T, q db.DB) (*auth.User, *auth.Session) {
	user, session := testauth.NewAuth(t, q)
	p := NewProfile()
	p.ID = ""
	p.User = user
	p.UserID = user.ID
	if err := p.Create(q); err != nil {
		t.Fatalf("failed to create a new auth with profile: %s", err)
	}
	return user, session
}

// NewAdminAuth creates a new admin and their session
func NewAdminAuth(t *testing.T, q db.DB) (*auth.User, *auth.Session) {
	user, session := testauth.NewAdminAuth(t, q)
	p := NewProfile()
	p.ID = ""
	p.User = user
	p.UserID = user.ID
	if err := p.Create(q); err != nil {
		t.Fatalf("failed to create a new admin auth with profile: %s", err)
	}
	return user, session
}
