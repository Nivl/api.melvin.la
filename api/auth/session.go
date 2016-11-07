package auth

import (
	"errors"
	"time"

	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/Nivl/api.melvin.la/api/db"
	uuid "github.com/satori/go.uuid"
)

// Session is a structure representing a session that can be saved in the database
type Session struct {
	ID        string     `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`

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

	s.ID = uuid.NewV4().String()
	s.CreatedAt = time.Now()

	stmt := "INSERT INTO sessions (id, created_at, user_id) VALUES ($1, $2, $3)"
	_, err := sql().Exec(stmt, s.ID, s.CreatedAt, s.UserID)
	return err
}

// FullyDelete deletes a session from the database
func (s *Session) FullyDelete() error {
	if s == nil {
		return errors.New("session not instanced")
	}

	if s.ID == "" {
		return errors.New("session has not been saved")
	}

	_, err := sql().Exec("DELETE FROM sessions WHERE id=$1", s.ID)
	return err
}

// Delete soft-deletes a session
func (s *Session) Delete() error {
	if s == nil {
		return apierror.NewServerError("session is not instanced")
	}

	if s.ID == "" {
		return apierror.NewServerError("cannot delete a non-persisted session")
	}

	now := time.Now()
	s.DeletedAt = &now

	stmt := `UPDATE sessions SET deleted_at = $2 WHERE id=$1`
	_, err := sql().Exec(stmt, s.ID, *s.DeletedAt)
	return err
}
