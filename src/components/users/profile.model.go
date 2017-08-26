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

	UserID           string `db:"user_id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Picture          string `db:"picture"`
	PhoneNumber      string `db:"phone_number"`
	PublicEmail      string `db:"public_email"`
	LinkedIn         string `db:"linkedin_custom_url"`
	FacebookUsername string `db:"facebook_username"`
	TwitterUsername  string `db:"twitter_username"`

	// Embedded models
	*auth.User `db:"user"`
}

// GetByIDWithProfile finds and returns an active user with their profile by ID
// Deleted object are not returned
func GetByIDWithProfile(q db.DB, id string) (*Profile, error) {
	u := &Profile{}
	stmt := `
	SELECT user.*, ` + JoinProfileSQL("profile") + `
	FROM users user
	JOIN user_profiles profile
	  ON user.id = profile.user_id
	WHERE user.id=$1
	  AND user.deleted_at IS NULL
	  AND profile.deleted_at IS NULL
	LIMIT 1`
	err := q.Get(u, stmt, id)
	return u, apierror.NewFromSQL(err)
}

// GetAnyByIDWithProfile finds and returns an active user with their
// profile by ID
// Deleted object are returned
func GetAnyByIDWithProfile(q db.DB, id string) (*Profile, error) {
	u := &Profile{}
	stmt := `
	SELECT user.*, ` + JoinProfileSQL("profile") + `
	FROM users user
	JOIN user_profiles profile
	  ON user.id = profile.user_id
	WHERE user.id=$1
	LIMIT 1`
	err := q.Get(u, stmt, id)
	return u, apierror.NewFromSQL(err)
}
