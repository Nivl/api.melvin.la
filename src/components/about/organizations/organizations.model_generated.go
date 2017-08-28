package organizations

// Code generated; DO NOT EDIT.

import (
	"errors"
	"fmt"

	"github.com/Nivl/go-rest-tools/types/apierror"
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
// Deleted object are not returned
func GetByID(q db.Queryable, id string) (*Organization, error) {
	o := &Organization{}
	stmt := "SELECT * from about_organizations WHERE id=$1 and deleted_at IS NULL LIMIT 1"
	err := q.Get(o, stmt, id)
	return o, apierror.NewFromSQL(err)
}

// GetAnyByID finds and returns an organization by ID.
// Deleted object are returned
func GetAnyByID(q db.Queryable, id string) (*Organization, error) {
	o := &Organization{}
	stmt := "SELECT * from about_organizations WHERE id=$1 LIMIT 1"
	err := q.Get(o, stmt, id)
	return o, apierror.NewFromSQL(err)
}


// Save creates or updates the article depending on the value of the id using
// a transaction
func (o *Organization) Save(q db.Queryable) error {
	if o.ID == "" {
		return o.Create(q)
	}

	return o.Update(q)
}

// Create persists a organization in the database
func (o *Organization) Create(q db.Queryable) error {
	if o.ID != "" {
		return errors.New("cannot persist a organization that already has an ID")
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

  return apierror.NewFromSQL(err)
}

// Update updates most of the fields of a persisted organization
// Excluded fields are id, created_at, deleted_at, etc.
func (o *Organization) Update(q db.Queryable) error {
	if o.ID == "" {
		return errors.New("cannot update a non-persisted organization")
	}

	return o.doUpdate(q)
}

// doUpdate updates a organization in the database
func (o *Organization) doUpdate(q db.Queryable) error {
	if o.ID == "" {
		return errors.New("cannot update a non-persisted organization")
	}

	o.UpdatedAt = db.Now()

	stmt := "UPDATE about_organizations SET id=:id, created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, name=:name, short_name=:short_name, logo=:logo, website=:website WHERE id=:id"
	_, err := q.NamedExec(stmt, o)

	return apierror.NewFromSQL(err)
}

// Delete removes a organization from the database
func (o *Organization) Delete(q db.Queryable) error {
	if o.ID == "" {
		return errors.New("organization has not been saved")
	}

	stmt := "DELETE FROM about_organizations WHERE id=$1"
	_, err := q.Exec(stmt, o.ID)

	return err
}

// GetID returns the ID field
func (o *Organization) GetID() string {
	return o.ID
}

// SetID sets the ID field
func (o *Organization) SetID(id string) {
	o.ID = id
}

// IsZero checks if the object is either nil or don't have an ID
func (o *Organization) IsZero() bool {
	return o == nil || o.ID == ""
}