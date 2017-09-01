package users

import (
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/types/apierror"
)

// Profile represents the public information of a user
//go:generate api-cli generate model Profile -t user_profiles -e Get,GetAny --single=false
type Profile struct {
	ID        string   `db:"id"`
	CreatedAt *db.Time `db:"created_at"`
	UpdatedAt *db.Time `db:"updated_at"`
	DeletedAt *db.Time `db:"deleted_at"`

	UserID           string  `db:"user_id"`
	FirstName        *string `json:"first_name"`
	LastName         *string `json:"last_name"`
	Picture          *string `db:"picture"`
	PhoneNumber      *string `db:"phone_number"`
	PublicEmail      *string `db:"public_email"`
	LinkedIn         *string `db:"linkedin_custom_url"`
	FacebookUsername *string `db:"facebook_username"`
	TwitterUsername  *string `db:"twitter_username"`
	IsFeatured       *bool   `db:"is_featured"`

	// Embedded models
	*auth.User `db:"users"`
}

// Profiles represents a list of Profile
type Profiles []*Profile

// GetByIDWithProfile finds and returns an active user with their profile by ID
// Deleted object are not returned
func GetByIDWithProfile(q db.Queryable, id string) (*Profile, error) {
	u := &Profile{}
	stmt := `
	SELECT profile.*, ` + auth.JoinUserSQL("users") + `
	FROM user_profiles profile
	JOIN users
	  ON users.id = profile.user_id
	WHERE users.id=$1
	  AND users.deleted_at IS NULL
	LIMIT 1`
	err := q.Get(u, stmt, id)
	return u, apierror.NewFromSQL(err)
}

// GetFeaturedProfile finds and returns the featured user
// Deleted object are not returned
func GetFeaturedProfile(q db.Queryable) (*Profile, error) {
	u := &Profile{}
	stmt := `
	SELECT profile.*, ` + auth.JoinUserSQL("users") + `
	FROM user_profiles profile
	JOIN users
	  ON users.id = profile.user_id
	WHERE users.is_featured IS true
	  AND users.deleted_at IS NULL
	LIMIT 1`
	err := q.Get(u, stmt)
	return u, apierror.NewFromSQL(err)
}
