package experience

// Code generated; DO NOT EDIT.

import (
	"errors"


	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-types/datetime"
	"github.com/Nivl/go-rest-tools/storage/db"
	uuid "github.com/satori/go.uuid"
)








// Save creates or updates the article depending on the value of the id using
// a transaction
func (e *Experience) Save(q db.Queryable) error {
	if e.ID == "" {
		return e.Create(q)
	}

	return e.Update(q)
}

// Create persists a experience in the database
func (e *Experience) Create(q db.Queryable) error {
	if e.ID != "" {
		return errors.New("cannot persist a experience that already has an ID")
	}

	return e.doCreate(q)
}

// doCreate persists a experience in the database using a Node
func (e *Experience) doCreate(q db.Queryable) error {
	e.ID = uuid.NewV4().String()
	e.UpdatedAt = datetime.Now()
	if e.CreatedAt == nil {
		e.CreatedAt = datetime.Now()
	}

	stmt := "INSERT INTO about_experience (id, created_at, updated_at, deleted_at, organization_id, job_title, location, description, start_date, end_date) VALUES (:id, :created_at, :updated_at, :deleted_at, :organization_id, :job_title, :location, :description, :start_date, :end_date)"
	_, err := q.NamedExec(stmt, e)

  return apierror.NewFromSQL(err)
}

// Update updates most of the fields of a persisted experience
// Excluded fields are id, created_at, deleted_at, etc.
func (e *Experience) Update(q db.Queryable) error {
	if e.ID == "" {
		return errors.New("cannot update a non-persisted experience")
	}

	return e.doUpdate(q)
}

// doUpdate updates a experience in the database
func (e *Experience) doUpdate(q db.Queryable) error {
	if e.ID == "" {
		return errors.New("cannot update a non-persisted experience")
	}

	e.UpdatedAt = datetime.Now()

	stmt := "UPDATE about_experience SET id=:id, created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, organization_id=:organization_id, job_title=:job_title, location=:location, description=:description, start_date=:start_date, end_date=:end_date WHERE id=:id"
	_, err := q.NamedExec(stmt, e)

	return apierror.NewFromSQL(err)
}

// Delete removes a experience from the database
func (e *Experience) Delete(q db.Queryable) error {
	if e.ID == "" {
		return errors.New("experience has not been saved")
	}

	stmt := "DELETE FROM about_experience WHERE id=$1"
	_, err := q.Exec(stmt, e.ID)

	return err
}

// GetID returns the ID field
func (e *Experience) GetID() string {
	return e.ID
}

// SetID sets the ID field
func (e *Experience) SetID(id string) {
	e.ID = id
}

// IsZero checks if the object is either nil or don't have an ID
func (e *Experience) IsZero() bool {
	return e == nil || e.ID == ""
}