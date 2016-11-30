package auth

// Code generated by api-cli; DO NOT EDIT\n

import (
	"errors"
	

	"github.com/jmoiron/sqlx"
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/app"
	"github.com/melvin-laplanche/ml-api/src/db"
	uuid "github.com/satori/go.uuid"
)



// Save creates or updates the session depending on the value of the id
func (s *Session) Save() error {
	return s.SaveTx(nil)
}



// Create persists a user in the database
func (s *Session) Create() error {
	return s.CreateTx(nil)
}



// doCreate persists an object in the database using an optional transaction
func (s *Session) doCreate(tx *sqlx.Tx) error {
	if s == nil {
		return errors.New("session not instanced")
	}

	s.ID = uuid.NewV4().String()
	s.CreatedAt = db.Now()
	s.UpdatedAt = db.Now()

	stmt := "INSERT INTO sessions (id, created_at, updated_at, deleted_at, user_id) VALUES (:id, :created_at, :updated_at, :deleted_at, :user_id)"
	var err error
	if tx == nil {
	  _, err = app.GetContext().SQL.NamedExec(stmt, s)
	} else {
		_, err = tx.NamedExec(stmt, s)
	}

  return err
}







// FullyDelete removes an object from the database
func (s *Session) FullyDelete() error {
	return s.FullyDeleteTx(nil)
}

// FullyDeleteTx removes an object from the database using a transaction
func (s *Session) FullyDeleteTx(tx *sqlx.Tx) error {
	if s == nil {
		return errors.New("session not instanced")
	}

	if s.ID == "" {
		return errors.New("session has not been saved")
	}

	stmt := "DELETE FROM sessions WHERE id=$1"
	var err error
	if tx == nil {
	  _, err = app.GetContext().SQL.Exec(stmt, s.ID)
	} else {
		_, err = tx.Exec(stmt, s.ID)
	}

	return err
}

// Delete soft delete an object.
func (s *Session) Delete() error {
	return s.DeleteTx(nil)
}

// DeleteTx soft delete an object using a transaction
func (s *Session) DeleteTx(tx *sqlx.Tx) error {
	return s.doDelete(tx)
}

// doDelete performs a soft delete operation on an object using an optional transaction
func (s *Session) doDelete(tx *sqlx.Tx) error {
	if s == nil {
		return apierror.NewServerError("session is not instanced")
	}

	if s.ID == "" {
		return apierror.NewServerError("cannot delete a non-persisted session")
	}

	s.DeletedAt = db.Now()

	stmt := "UPDATE sessions SET deleted_at = $2 WHERE id=$1"
	var err error
	if tx == nil {
	  _, err = app.GetContext().SQL.Exec(stmt, s.ID, s.DeletedAt)
	} else {
		_, err = tx.Exec(stmt, s.ID, s.DeletedAt)
	}
	return err
}

// IsZero checks if the object is either nil or don't have an ID
func (s *Session) IsZero() bool {
	return s == nil || s.ID == ""
}