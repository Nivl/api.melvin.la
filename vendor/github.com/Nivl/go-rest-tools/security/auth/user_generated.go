package auth

// Code auto-generated; DO NOT EDIT

import (
	"errors"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/storage/db"
	uuid "github.com/satori/go.uuid"
)

// UserExists checks if a user exists for a specific ID
func UserExists(id string) (bool, error) {
	exists := false
	stmt := "SELECT exists(SELECT 1 FROM users WHERE id=$1 and deleted_at IS NULL)"
	err := db.Writer.Get(&exists, stmt, id)
	return exists, err
}

// Save creates or updates the user depending on the value of the id
func (u *User) Save() error {
	return u.SaveQ(db.Writer)
}

// SaveQ creates or updates the article depending on the value of the id using
// a transaction
func (u *User) SaveQ(q db.Queryable) error {
	if u == nil {
		return httperr.NewServerError("user is not instanced")
	}

	if u.ID == "" {
		return u.CreateQ(q)
	}

	return u.UpdateQ(q)
}

// Create persists a user in the database
func (u *User) Create() error {
	return u.CreateQ(db.Writer)
}

// doCreate persists a user in the database using a Node
func (u *User) doCreate(q db.Queryable) error {
	if u == nil {
		return errors.New("user not instanced")
	}

	u.ID = uuid.NewV4().String()
	u.CreatedAt = db.Now()
	u.UpdatedAt = db.Now()

	stmt := "INSERT INTO users (id, created_at, updated_at, deleted_at, name, email, password, is_admin) VALUES (:id, :created_at, :updated_at, :deleted_at, :name, :email, :password, :is_admin)"
	_, err := q.NamedExec(stmt, u)

	return err
}

// Update updates most of the fields of a persisted user.
// Excluded fields are id, created_at, deleted_at, etc.
func (u *User) Update() error {
	return u.UpdateQ(db.Writer)
}

// doUpdate updates a user in the database using an optional transaction
func (u *User) doUpdate(q db.Queryable) error {
	if u == nil {
		return httperr.NewServerError("user is not instanced")
	}

	if u.ID == "" {
		return httperr.NewServerError("cannot update a non-persisted user")
	}

	u.UpdatedAt = db.Now()

	stmt := "UPDATE users SET id=:id, created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, name=:name, email=:email, password=:password, is_admin=:is_admin WHERE id=:id"
	_, err := q.NamedExec(stmt, u)

	return err
}

// FullyDelete removes a user from the database
func (u *User) FullyDelete() error {
	return u.FullyDeleteQ(db.Writer)
}

// FullyDeleteQ removes a user from the database using a transaction
func (u *User) FullyDeleteQ(q db.Queryable) error {
	if u == nil {
		return errors.New("user not instanced")
	}

	if u.ID == "" {
		return errors.New("user has not been saved")
	}

	stmt := "DELETE FROM users WHERE id=$1"
	_, err := q.Exec(stmt, u.ID)

	return err
}

// Delete soft delete a user.
func (u *User) Delete() error {
	return u.DeleteQ(db.Writer)
}

// DeleteQ soft delete a user using a transaction
func (u *User) DeleteQ(q db.Queryable) error {
	return u.doDelete(q)
}

// doDelete performs a soft delete operation on a user using an optional transaction
func (u *User) doDelete(q db.Queryable) error {
	if u == nil {
		return httperr.NewServerError("user is not instanced")
	}

	if u.ID == "" {
		return httperr.NewServerError("cannot delete a non-persisted user")
	}

	u.DeletedAt = db.Now()

	stmt := "UPDATE users SET deleted_at = $2 WHERE id=$1"
	_, err := q.Exec(stmt, u.ID, u.DeletedAt)
	return err
}

// IsZero checks if the object is either nil or don't have an ID
func (u *User) IsZero() bool {
	return u == nil || u.ID == ""
}
