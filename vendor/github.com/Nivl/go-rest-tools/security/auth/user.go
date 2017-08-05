package auth

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/Nivl/go-rest-tools/storage/db"
)

// User is a structure representing a user that can be saved in the database
//go:generate api-cli generate model User -t users --single=false
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
