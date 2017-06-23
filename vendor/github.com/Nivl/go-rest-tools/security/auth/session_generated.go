package auth

// Code auto-generated; DO NOT EDIT

import (
	"errors"
	

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/storage/db"
	uuid "github.com/satori/go.uuid"
	"github.com/jmoiron/sqlx"
)











// doCreate persists a session in the database using a Node
func (s *Session) doCreate(q *sqlx.DB) error {
	if s == nil {
		return errors.New("session not instanced")
	}

	s.ID = uuid.NewV4().String()
	s.UpdatedAt = db.Now()
	if s.CreatedAt == nil {
		s.CreatedAt = db.Now()
	}

	stmt := "INSERT INTO sessions (id, created_at, updated_at, deleted_at, user_id) VALUES (:id, :created_at, :updated_at, :deleted_at, :user_id)"
	_, err := q.NamedExec(stmt, s)

  return err
}





// Delete removes a session from the database using a transaction
func (s *Session) Delete(q *sqlx.DB) error {
	if s == nil {
		return errors.New("session not instanced")
	}

	if s.ID == "" {
		return errors.New("session has not been saved")
	}

	stmt := "DELETE FROM sessions WHERE id=$1"
	_, err := q.Exec(stmt, s.ID)

	return err
}

// Trash soft delete a session using a transaction
func (s *Session) Trash(q *sqlx.DB) error {
	return s.doTrash(q)
}

// doTrash performs a soft delete operation on a session using an optional transaction
func (s *Session) doTrash(q *sqlx.DB) error {
	if s == nil {
		return httperr.NewServerError("session is not instanced")
	}

	if s.ID == "" {
		return httperr.NewServerError("cannot trash a non-persisted session")
	}

	s.DeletedAt = db.Now()

	stmt := "UPDATE sessions SET deleted_at = $2 WHERE id=$1"
	_, err := q.Exec(stmt, s.ID, s.DeletedAt)
	return err
}

// IsZero checks if the object is either nil or don't have an ID
func (s *Session) IsZero() bool {
	return s == nil || s.ID == ""
}