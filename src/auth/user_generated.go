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



// Save creates or updates the user depending on the value of the id
func (u *User) Save() error {
	return u.SaveTx(nil)
}

// SaveTx creates or updates the article depending on the value of the id using
// a transaction
func (u *User) SaveTx(tx *sqlx.Tx) error {
	if u == nil {
		return apierror.NewServerError("user is not instanced")
	}

	if u.ID == "" {
		return u.CreateTx(tx)
	}

	return u.UpdateTx(tx)
}

// Create persists a user in the database
func (u *User) Create() error {
	return u.CreateTx(nil)
}



// doCreate persists an object in the database using an optional transaction
func (u *User) doCreate(tx *sqlx.Tx) error {
	if u == nil {
		return errors.New("user not instanced")
	}

	u.ID = uuid.NewV4().String()
	u.CreatedAt = db.Now()
	u.UpdatedAt = db.Now()

	stmt := "INSERT INTO users (id, created_at, updated_at, deleted_at, name, email, password) VALUES (:id, :created_at, :updated_at, :deleted_at, :name, :email, :password)"
	var err error
	if tx == nil {
	  _, err = app.GetContext().SQL.NamedExec(stmt, u)
	} else {
		_, err = tx.NamedExec(stmt, u)
	}

  return err
}

// Update updates most of the fields of a persisted user.
// Excluded fields are id, created_at, deleted_at, etc.
func (u *User) Update() error {
	return u.UpdateTx(nil)
}



// doUpdate updates an object in the database using an optional transaction
func (u *User) doUpdate(tx *sqlx.Tx) error {
	if u == nil {
		return apierror.NewServerError("user is not instanced")
	}

	if u.ID == "" {
		return apierror.NewServerError("cannot update a non-persisted user")
	}

	u.UpdatedAt = db.Now()

	stmt := "UPDATE users SET id=:id, created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, name=:name, email=:email, password=:password WHERE id=:id"
	var err error
	if tx == nil {
	  _, err = app.GetContext().SQL.NamedExec(stmt, u)
	} else {
		_, err = tx.NamedExec(stmt, u)
	}

	return err
}

// FullyDelete removes an object from the database
func (u *User) FullyDelete() error {
	return u.FullyDeleteTx(nil)
}

// FullyDeleteTx removes an object from the database using a transaction
func (u *User) FullyDeleteTx(tx *sqlx.Tx) error {
	if u == nil {
		return errors.New("user not instanced")
	}

	if u.ID == "" {
		return errors.New("user has not been saved")
	}

	stmt := "DELETE FROM users WHERE id=$1"
	var err error
	if tx == nil {
	  _, err = app.GetContext().SQL.Exec(stmt, u.ID)
	} else {
		_, err = tx.Exec(stmt, u.ID)
	}

	return err
}

// Delete soft delete an object.
func (u *User) Delete() error {
	return u.DeleteTx(nil)
}

// DeleteTx soft delete an object using a transaction
func (u *User) DeleteTx(tx *sqlx.Tx) error {
	return u.doDelete(tx)
}

// doDelete performs a soft delete operation on an object using an optional transaction
func (u *User) doDelete(tx *sqlx.Tx) error {
	if u == nil {
		return apierror.NewServerError("user is not instanced")
	}

	if u.ID == "" {
		return apierror.NewServerError("cannot delete a non-persisted user")
	}

	u.DeletedAt = db.Now()

	stmt := "UPDATE users SET deleted_at = $2 WHERE id=$1"
	var err error
	if tx == nil {
	  _, err = app.GetContext().SQL.Exec(stmt, u.ID, u.DeletedAt)
	} else {
		_, err = tx.Exec(stmt, u.ID, u.DeletedAt)
	}
	return err
}

// IsZero checks if the object is either nil or don't have an ID
func (u *User) IsZero() bool {
	return u == nil || u.ID == ""
}