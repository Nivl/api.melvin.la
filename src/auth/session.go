package auth

import (
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/db"
)

// Session is a structure representing a session that can be saved in the database
//go:generate api-cli generate model Session -t sessions
type Session struct {
	ID        string   `db:"id"`
	CreatedAt *db.Time `db:"created_at"`
	DeletedAt *db.Time `db:"deleted_at"`

	UserID string `db:"user_id"`
}

// Exists check if a session exists in the database
func (s *Session) Exists() (bool, error) {
	if s == nil {
		return false, apierror.NewServerError("session is nil")
	}

	if s.UserID == "" {
		return false, apierror.NewServerError("user id required")
	}

	// Deleted sessions should be explicitly checked
	if s.DeletedAt != nil {
		return false, nil
	}

	var count int
	stmt := `SELECT count(1)
					FROM sessions
					WHERE deleted_at IS NULL
						AND id = $1
						AND user_id = $2`
	err := db.Get(&count, stmt, s.ID, s.UserID)
	return (count > 0), err
}

// Save is an alias for Create since sessions are not updatable
func (s *Session) Save() error {
	if s == nil {
		return apierror.NewServerError("session is nil")
	}

	return s.Create()
}

// Create persists a session in the database
func (s *Session) Create() error {
	if s == nil {
		return apierror.NewServerError("session is nil")
	}

	if s.ID != "" {
		return apierror.NewServerError("sessions cannot be updated")
	}

	if s.UserID == "" {
		return apierror.NewServerError("cannot save a session with no user id")
	}

	return s.doCreate()
}

// Delete soft-deletes a session
func (s *Session) Delete() error {
	return s.doCreate()
}
