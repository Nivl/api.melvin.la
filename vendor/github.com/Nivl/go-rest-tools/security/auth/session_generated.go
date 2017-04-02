package auth

// Code auto-generated; DO NOT EDIT

import (
	"errors"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/storage/db"
	uuid "github.com/satori/go.uuid"
)

// Save creates or updates the session depending on the value of the id
func (s *Session) Save() error {
	return s.SaveQ(db.Writer)
}

// Create persists a session in the database
func (s *Session) Create() error {
	return s.CreateQ(db.Writer)
}

// doCreate persists a session in the database using a Node
func (s *Session) doCreate(q db.Queryable) error {
	if s == nil {
		return errors.New("session not instanced")
	}

	s.ID = uuid.NewV4().String()
	s.CreatedAt = db.Now()
	s.UpdatedAt = db.Now()

	stmt := "INSERT INTO sessions (id, created_at, updated_at, deleted_at, user_id) VALUES (:id, :created_at, :updated_at, :deleted_at, :user_id)"
	_, err := q.NamedExec(stmt, s)

	return err
}

// FullyDelete removes a session from the database
func (s *Session) FullyDelete() error {
	return s.FullyDeleteQ(db.Writer)
}

// FullyDeleteQ removes a session from the database using a transaction
func (s *Session) FullyDeleteQ(q db.Queryable) error {
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

// Delete soft delete a session.
func (s *Session) Delete() error {
	return s.DeleteQ(db.Writer)
}

// DeleteQ soft delete a session using a transaction
func (s *Session) DeleteQ(q db.Queryable) error {
	return s.doDelete(q)
}

// doDelete performs a soft delete operation on a session using an optional transaction
func (s *Session) doDelete(q db.Queryable) error {
	if s == nil {
		return httperr.NewServerError("session is not instanced")
	}

	if s.ID == "" {
		return httperr.NewServerError("cannot delete a non-persisted session")
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
