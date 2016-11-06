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
	UUID      string     `db:"uuid"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`

	UserUUID string `db:"user_uuid"`
}

// Exists check if a session exists in the database
func (s *Session) Exists() (bool, error) {
	if s == nil {
		return false, apierror.NewServerError("session is nil")
	}

	if s.UserUUID == "" {
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
						AND uuid = $1
						AND user_uuid = $2`
	err := db.Get(&count, stmt, s.UUID, s.UserUUID)
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

	if s.UUID != "" {
		return apierror.NewServerError("sessions cannot be updated")
	}

	if s.UserUUID == "" {
		return apierror.NewServerError("cannot save a session with no user uuid")
	}

	s.UUID = uuid.NewV4().String()
	s.CreatedAt = time.Now()

	stmt := "INSERT INTO sessions (uuid, created_at, user_uuid) VALUES ($1, $2, $3)"
	_, err := sql().Exec(stmt, s.UUID, s.CreatedAt, s.UserUUID)
	return err
}

// FullyDelete deletes a session from the database
func (s *Session) FullyDelete() error {
	if s == nil {
		return errors.New("session not instanced")
	}

	if s.UUID == "" {
		return errors.New("session has not been saved")
	}

	_, err := sql().Exec("DELETE FROM sessions WHERE uuid=$1", s.UUID)
	return err
}

// Delete soft-deletes a session
func (s *Session) Delete() error {
	if s == nil {
		return apierror.NewServerError("session is not instanced")
	}

	if s.UUID == "" {
		return apierror.NewServerError("cannot delete a non-persisted session")
	}

	now := time.Now()
	s.DeletedAt = &now

	stmt := `UPDATE sessions SET deleted_at = $2 WHERE uuid=$1`
	_, err := sql().Exec(stmt, s.UUID, *s.DeletedAt)
	return err
}
