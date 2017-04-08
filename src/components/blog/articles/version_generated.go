package articles

// Code auto-generated; DO NOT EDIT

import (
	"errors"
	"fmt"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/storage/db"
	uuid "github.com/satori/go.uuid"
)

// JoinVersionSQL returns a string ready to be embed in a JOIN query
func JoinVersionSQL(prefix string) string {
	fields := []string{ "id", "created_at", "updated_at", "deleted_at", "article_id", "title", "content", "subtitle", "description" }
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

// GetVersionByID finds and returns an active version by ID
func GetVersionByID(id string) (*Version, error) {
	v := &Version{}
	stmt := "SELECT * from blog_article_versions WHERE id=$1 and deleted_at IS NULL LIMIT 1"
	err := db.Get(v, stmt, id)
	// We want to return nil if a version is not found
	if v.ID == "" {
		return nil, err
	}
	return v, err
}

// VersionExists checks if a version exists for a specific ID
func VersionExists(id string) (bool, error) {
	exists := false
	stmt := "SELECT exists(SELECT 1 FROM blog_article_versions WHERE id=$1 and deleted_at IS NULL)"
	err := db.Writer.Get(&exists, stmt, id)
	return exists, err
}

// Save creates or updates the version depending on the value of the id
func (v *Version) Save() error {
	return v.SaveQ(db.Writer)
}

// SaveQ creates or updates the article depending on the value of the id using
// a transaction
func (v *Version) SaveQ(q db.Queryable) error {
	if v == nil {
		return httperr.NewServerError("version is not instanced")
	}

	if v.ID == "" {
		return v.CreateQ(q)
	}

	return v.UpdateQ(q)
}

// Create persists a version in the database
func (v *Version) Create() error {
	return v.CreateQ(db.Writer)
}

// CreateQ persists a version in the database
func (v *Version) CreateQ(q db.Queryable) error {
	if v == nil {
		return httperr.NewServerError("version is not instanced")
	}

	if v.ID != "" {
		return httperr.NewServerError("cannot persist a version that already has an ID")
	}

	return v.doCreate(q)
}

// doCreate persists a version in the database using a Node
func (v *Version) doCreate(q db.Queryable) error {
	if v == nil {
		return errors.New("version not instanced")
	}

	v.ID = uuid.NewV4().String()
	v.UpdatedAt = db.Now()
	if v.CreatedAt == nil {
		v.CreatedAt = db.Now()
	}

	stmt := "INSERT INTO blog_article_versions (id, created_at, updated_at, deleted_at, article_id, title, content, subtitle, description) VALUES (:id, :created_at, :updated_at, :deleted_at, :article_id, :title, :content, :subtitle, :description)"
	_, err := q.NamedExec(stmt, v)

  return err
}

// Update updates most of the fields of a persisted version.
// Excluded fields are id, created_at, deleted_at, etc.
func (v *Version) Update() error {
	return v.UpdateQ(db.Writer)
}

// UpdateQ updates most of the fields of a persisted version using a transaction
// Excluded fields are id, created_at, deleted_at, etc.
func (v *Version) UpdateQ(q db.Queryable) error {
	if v == nil {
		return httperr.NewServerError("version is not instanced")
	}

	if v.ID == "" {
		return httperr.NewServerError("cannot update a non-persisted version")
	}

	return v.doUpdate(q)
}

// doUpdate updates a version in the database using an optional transaction
func (v *Version) doUpdate(q db.Queryable) error {
	if v == nil {
		return httperr.NewServerError("version is not instanced")
	}

	if v.ID == "" {
		return httperr.NewServerError("cannot update a non-persisted version")
	}

	v.UpdatedAt = db.Now()

	stmt := "UPDATE blog_article_versions SET id=:id, created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, article_id=:article_id, title=:title, content=:content, subtitle=:subtitle, description=:description WHERE id=:id"
	_, err := q.NamedExec(stmt, v)

	return err
}

// FullyDelete removes a version from the database
func (v *Version) FullyDelete() error {
	return v.FullyDeleteQ(db.Writer)
}

// FullyDeleteQ removes a version from the database using a transaction
func (v *Version) FullyDeleteQ(q db.Queryable) error {
	if v == nil {
		return errors.New("version not instanced")
	}

	if v.ID == "" {
		return errors.New("version has not been saved")
	}

	stmt := "DELETE FROM blog_article_versions WHERE id=$1"
	_, err := q.Exec(stmt, v.ID)

	return err
}

// Delete soft delete a version.
func (v *Version) Delete() error {
	return v.DeleteQ(db.Writer)
}

// DeleteQ soft delete a version using a transaction
func (v *Version) DeleteQ(q db.Queryable) error {
	return v.doDelete(q)
}

// doDelete performs a soft delete operation on a version using an optional transaction
func (v *Version) doDelete(q db.Queryable) error {
	if v == nil {
		return httperr.NewServerError("version is not instanced")
	}

	if v.ID == "" {
		return httperr.NewServerError("cannot delete a non-persisted version")
	}

	v.DeletedAt = db.Now()

	stmt := "UPDATE blog_article_versions SET deleted_at = $2 WHERE id=$1"
	_, err := q.Exec(stmt, v.ID, v.DeletedAt)
	return err
}

// IsZero checks if the object is either nil or don't have an ID
func (v *Version) IsZero() bool {
	return v == nil || v.ID == ""
}