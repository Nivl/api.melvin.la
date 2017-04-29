package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/storage/db"
)

// User is a structure representing a user that can be saved in the database
//go:generate tv-api-cli generate model User -t users -e CreateQ,UpdateQ,Get,JoinSQL
type User struct {
	ID        string   `db:"id"`
	CreatedAt *db.Time `db:"created_at"`
	UpdatedAt *db.Time `db:"updated_at"`
	DeletedAt *db.Time `db:"deleted_at"`

	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
	IsAdmin  bool   `db:"is_admin"`
}

// UserJoinSQL returns a string ready to be embed in a JOIN query
func UserJoinSQL(prefix string) string {
	fields := []string{"id", "created_at", "updated_at", "deleted_at", "name", "email", "password"}
	output := ""

	for i, field := range fields {
		if i != 0 {
			output += ", "
		}

		fullName := fmt.Sprintf("%s.%s", prefix, field)
		output += fmt.Sprintf("%s \"%s\"", fullName, fullName)
	}

	return output
}

// GetUser finds and returns an active user by ID
func GetUser(id string) (*User, error) {
	user := &User{}
	stmt := "SELECT * from users WHERE id=$1 and deleted_at IS NULL LIMIT 1"
	err := db.Get(user, stmt, id)
	// We want to return nil if a user is not found
	if user.ID == "" {
		return nil, err
	}
	return user, err
}

// CryptPassword returns a password encrypted with bcrypt
func CryptPassword(raw string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(password), nil
}

// IsPasswordValid Compare a bcrypt hash with a raw string and check if they match
func IsPasswordValid(hash string, raw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw))
	return err == nil
}

// CreateQ persists a user in the database
func (u *User) CreateQ(q db.Queryable) error {
	if u == nil {
		return httperr.NewServerError("user is not instanced")
	}

	if u.ID != "" {
		return httperr.NewServerError("cannot persist a user that already has a ID")
	}

	err := u.doCreate(q)
	if err != nil && db.IsDup(err) {
		return httperr.NewConflict("email address already in use")
	}

	return err
}

// UpdateQ updates most of the fields of a persisted user.
// Excluded fields are id, created_at, deleted_at
func (u *User) UpdateQ(q db.Queryable) error {
	if u == nil {
		return httperr.NewServerError("user is not instanced")
	}

	if u.ID == "" {
		return httperr.NewServerError("cannot update a non-persisted user")
	}

	err := u.doUpdate(q)
	if err != nil && db.IsDup(err) {
		return httperr.NewConflict("email address already in use")
	}

	return err
}

// IsLogged checks if the user object belong to a logged in user
// Works on nil object
func (u *User) IsLogged() bool {
	return u != nil
}

// IsAdm checks if the user object belong to a logged in admin
// Works on nil object
func (u *User) IsAdm() bool {
	return u.IsLogged() && u.IsAdmin
}
