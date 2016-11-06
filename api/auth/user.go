package auth

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/Nivl/api.melvin.la/api/db"
	"github.com/dchest/uniuri"
	uuid "github.com/satori/go.uuid"
)

// User is a structure representing a user that can be saved in the database
type User struct {
	UUID      string     `db:"uuid"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`

	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

// GetUser finds and returns an active user by UUID
func GetUser(uuid string) (*User, error) {
	user := &User{}
	stmt := "SELECT * from users WHERE uuid=$1 and deleted_at IS NULL"
	err := db.Get(user, stmt, uuid)
	// We want to return nil if a user is not found
	if user.UUID == "" {
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

// Save creates or updates the user depending on the value of the uuid
func (u *User) Save() error {
	if u == nil {
		return apierror.NewServerError("user is not instanced")
	}

	if u.UUID == "" {
		return u.Create()
	}

	return u.Update()
}

// FullyDelete remove the user from the database
func (u *User) FullyDelete() error {
	if u == nil {
		return errors.New("user not instanced")
	}

	if u.UUID == "" {
		return errors.New("user has not been saved")
	}

	_, err := sql().Exec("DELETE FROM users WHERE uuid=$1", u.UUID)
	return err
}

// Create persists a user in the database
func (u *User) Create() error {
	if u == nil {
		return apierror.NewServerError("user is not instanced")
	}

	if u.UUID != "" {
		return apierror.NewServerError("cannot persist a user that already has a UUID")
	}

	u.UUID = uuid.NewV4().String()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	stmt := "INSERT INTO users (uuid, created_at, updated_at, name, email, password) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := sql().Exec(stmt, u.UUID, u.CreatedAt, u.UpdatedAt, u.Name, u.Email, u.Password)

	if err != nil && db.SQLIsDup(err) {
		return apierror.NewConflict("email address already in use")
	}

	return err
}

// Update updates most of the fields of a persisted user.
// Excluded fields are uuid, created_at, deleted_at
func (u *User) Update() error {
	if u == nil {
		return apierror.NewServerError("user is not instanced")
	}

	if u.UUID == "" {
		return apierror.NewServerError("cannot update a non-persisted user")
	}

	u.UpdatedAt = time.Now()

	stmt := `UPDATE users
					 SET updated_at = $2,
					 		name = $3,
							email = $4,
							password = $5
	         WHERE uuid=$1`
	_, err := sql().Exec(stmt, u.UUID, u.UpdatedAt, u.Name, u.Email, u.Password)

	if err != nil && db.SQLIsDup(err) {
		return apierror.NewConflict("email address already in use")
	}

	return err
}

// Delete soft delete a user.
func (u *User) Delete() error {
	if u == nil {
		return apierror.NewServerError("user is not instanced")
	}

	if u.UUID == "" {
		return apierror.NewServerError("cannot delete a non-persisted user")
	}

	now := time.Now()
	u.DeletedAt = &now

	stmt := `UPDATE users SET deleted_at = $2 WHERE uuid=$1`
	_, err := sql().Exec(stmt, u.UUID, *u.DeletedAt)
	return err
}

// NewTestUser creates a new user with "fake" as password
func NewTestUser(t *testing.T, u *User) *User {
	if u == nil {
		u = &User{}
	}

	if u.Email == "" {
		u.Email = fmt.Sprintf("fake+%s@melvin.la", uniuri.New())
	}

	if u.Name == "" {
		u.Name = "Fake Account"
	}

	if u.Password == "" {
		var err error
		u.Password, err = CryptPassword("fake")
		if err != nil {
			t.Fatalf("failed to create password: %s", err)
		}
	}

	if err := u.Create(); err != nil {
		t.Fatalf("failed to create user: %s", err)
	}
	return u
}
