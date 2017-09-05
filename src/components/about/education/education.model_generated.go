package education

// Code generated; DO NOT EDIT.

import (
	"errors"
	

	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-rest-tools/types/datetime"
	"github.com/Nivl/go-rest-tools/storage/db"
	uuid "github.com/satori/go.uuid"
)








// Save creates or updates the article depending on the value of the id using
// a transaction
func (e *Education) Save(q db.Queryable) error {
	if e.ID == "" {
		return e.Create(q)
	}

	return e.Update(q)
}

// Create persists a education in the database
func (e *Education) Create(q db.Queryable) error {
	if e.ID != "" {
		return errors.New("cannot persist a education that already has an ID")
	}

	return e.doCreate(q)
}

// doCreate persists a education in the database using a Node
func (e *Education) doCreate(q db.Queryable) error {
	e.ID = uuid.NewV4().String()
	e.UpdatedAt = datetime.Now()
	if e.CreatedAt == nil {
		e.CreatedAt = datetime.Now()
	}

	stmt := "INSERT INTO about_education (id, created_at, updated_at, deleted_at, organization_id, degree, gpa, location, description, start_year, end_year) VALUES (:id, :created_at, :updated_at, :deleted_at, :organization_id, :degree, :gpa, :location, :description, :start_year, :end_year)"
	_, err := q.NamedExec(stmt, e)

  return apierror.NewFromSQL(err)
}

// Update updates most of the fields of a persisted education
// Excluded fields are id, created_at, deleted_at, etc.
func (e *Education) Update(q db.Queryable) error {
	if e.ID == "" {
		return errors.New("cannot update a non-persisted education")
	}

	return e.doUpdate(q)
}

// doUpdate updates a education in the database
func (e *Education) doUpdate(q db.Queryable) error {
	if e.ID == "" {
		return errors.New("cannot update a non-persisted education")
	}

	e.UpdatedAt = datetime.Now()

	stmt := "UPDATE about_education SET id=:id, created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, organization_id=:organization_id, degree=:degree, gpa=:gpa, location=:location, description=:description, start_year=:start_year, end_year=:end_year WHERE id=:id"
	_, err := q.NamedExec(stmt, e)

	return apierror.NewFromSQL(err)
}

// Delete removes a education from the database
func (e *Education) Delete(q db.Queryable) error {
	if e.ID == "" {
		return errors.New("education has not been saved")
	}

	stmt := "DELETE FROM about_education WHERE id=$1"
	_, err := q.Exec(stmt, e.ID)

	return err
}

// GetID returns the ID field
func (e *Education) GetID() string {
	return e.ID
}

// SetID sets the ID field
func (e *Education) SetID(id string) {
	e.ID = id
}

// IsZero checks if the object is either nil or don't have an ID
func (e *Education) IsZero() bool {
	return e == nil || e.ID == ""
}