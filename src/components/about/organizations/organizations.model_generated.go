package organizations

// Code auto-generated; DO NOT EDIT

import (
	"errors"
	"fmt"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/storage/db"
	uuid "github.com/satori/go.uuid"
)

// JoinSQL returns a string ready to be embed in a JOIN query
func JoinSQL(prefix string) string {
	fields := []string{ "id", "created_at", "updated_at", "deleted_at", "name", "short_name", "logo", "website" }
	output := ""

	for i, field := range fields {
		if i != 0 {
			output += ", "
		}

		fullName := fmt.Sprintf("%s.%s", prefix, field)
		output += fmt.Sprintf("%s \"%s\"", fullName, fullName)
	}

	return output
}

// GetByID finds and returns an active organization by ID
func GetByID(id string) (*Organization, error) {
	o := &Organization{}
	stmt := "SELECT * from about_organizations WHERE id=$1 and deleted_at IS NULL LIMIT 1"
	err := db.Get(o, stmt, id)
	// We want to return nil if a organization is not found
	if o.ID == "" {
		return nil, err
	}
	return o, err
}

// Exists checks if a organization exists for a specific ID
func Exists(id string) (bool, error) {
	exists := false
	stmt := "SELECT exists(SELECT 1 FROM about_organizations WHERE id=$1 and deleted_at IS NULL)"
	err := db.Writer.Get(&exists, stmt, id)
	return exists, err
}

// Save creates or updates the organization depending on the value of the id
func (o *Organization) Save() error {
	return o.SaveQ(db.Writer)
}

// SaveQ creates or updates the article depending on the value of the id using
// a transaction
func (o *Organization) SaveQ(q db.Queryable) error {
	if o == nil {
		return httperr.NewServerError("organization is not instanced")
	}

	if o.ID == "" {
		return o.CreateQ(q)
	}

	return o.UpdateQ(q)
}

// Create persists a organization in the database
func (o *Organization) Create() error {
	return o.CreateQ(db.Writer)
}

// CreateQ persists a organization in the database
func (o *Organization) CreateQ(q db.Queryable) error {
	if o == nil {
		return httperr.NewServerError("organization is not instanced")
	}

	if o.ID != "" {
		return httperr.NewServerError("cannot persist a organization that already has an ID")
	}

	return o.doCreate(q)
}

// doCreate persists a organization in the database using a Node
func (o *Organization) doCreate(q db.Queryable) error {
	if o == nil {
		return errors.New("organization not instanced")
	}

	o.ID = uuid.NewV4().String()
	o.UpdatedAt = db.Now()
	if o.CreatedAt == nil {
		o.CreatedAt = db.Now()
	}

	stmt := "INSERT INTO about_organizations (id, created_at, updated_at, deleted_at, name, short_name, logo, website) VALUES (:id, :created_at, :updated_at, :deleted_at, :name, :short_name, :logo, :website)"
	_, err := q.NamedExec(stmt, o)

  return err
}

// Update updates most of the fields of a persisted organization.
// Excluded fields are id, created_at, deleted_at, etc.
func (o *Organization) Update() error {
	return o.UpdateQ(db.Writer)
}

// UpdateQ updates most of the fields of a persisted organization using a transaction
// Excluded fields are id, created_at, deleted_at, etc.
func (o *Organization) UpdateQ(q db.Queryable) error {
	if o == nil {
		return httperr.NewServerError("organization is not instanced")
	}

	if o.ID == "" {
		return httperr.NewServerError("cannot update a non-persisted organization")
	}

	return o.doUpdate(q)
}

// doUpdate updates a organization in the database using an optional transaction
func (o *Organization) doUpdate(q db.Queryable) error {
	if o == nil {
		return httperr.NewServerError("organization is not instanced")
	}

	if o.ID == "" {
		return httperr.NewServerError("cannot update a non-persisted organization")
	}

	o.UpdatedAt = db.Now()

	stmt := "UPDATE about_organizations SET id=:id, created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, name=:name, short_name=:short_name, logo=:logo, website=:website WHERE id=:id"
	_, err := q.NamedExec(stmt, o)

	return err
}

// FullyDelete removes a organization from the database
func (o *Organization) FullyDelete() error {
	return o.FullyDeleteQ(db.Writer)
}

// FullyDeleteQ removes a organization from the database using a transaction
func (o *Organization) FullyDeleteQ(q db.Queryable) error {
	if o == nil {
		return errors.New("organization not instanced")
	}

	if o.ID == "" {
		return errors.New("organization has not been saved")
	}

	stmt := "DELETE FROM about_organizations WHERE id=$1"
	_, err := q.Exec(stmt, o.ID)

	return err
}

// Delete soft delete a organization.
func (o *Organization) Delete() error {
	return o.DeleteQ(db.Writer)
}

// DeleteQ soft delete a organization using a transaction
func (o *Organization) DeleteQ(q db.Queryable) error {
	return o.doDelete(q)
}

// doDelete performs a soft delete operation on a organization using an optional transaction
func (o *Organization) doDelete(q db.Queryable) error {
	if o == nil {
		return httperr.NewServerError("organization is not instanced")
	}

	if o.ID == "" {
		return httperr.NewServerError("cannot delete a non-persisted organization")
	}

	o.DeletedAt = db.Now()

	stmt := "UPDATE about_organizations SET deleted_at = $2 WHERE id=$1"
	_, err := q.Exec(stmt, o.ID, o.DeletedAt)
	return err
}

// IsZero checks if the object is either nil or don't have an ID
func (o *Organization) IsZero() bool {
	return o == nil || o.ID == ""
}