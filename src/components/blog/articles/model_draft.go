package articles

import (
	"errors"

	"github.com/Nivl/sqalx"
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/db"
)

// Draft Represents an article draft
type Draft Content

// ToContent returns a Content from a Draft
func (d *Draft) ToContent() *Content {
	content := Content(*d)
	return &content
}

// CreateTx persists a draft in the database using a transaction
func (d *Draft) CreateTx(tx sqalx.Node) error {
	return d.ToContent().CreateTx(tx)
}

// Create persists a draft in the database
func (d *Draft) Create() error {
	return d.ToContent().Create()
}

// UpdateTx updates most of the fields of a persisted content using a transaction
// Excluded fields are id, created_at, deleted_at
func (d *Draft) UpdateTx(tx sqalx.Node) error {
	return d.ToContent().UpdateTx(tx)
}

// Update updates most of the fields of a persisted content
// Excluded fields are id, created_at, deleted_at
func (d *Draft) Update() error {
	return d.ToContent().Update()
}

// SaveTx creates or updates the content depending on the value of the id using a transaction
func (d *Draft) SaveTx(tx sqalx.Node) error {
	return d.ToContent().SaveTx(tx)
}

// Save creates or updates the content depending on the value of the id
func (d *Draft) Save() error {
	return d.ToContent().Save()
}

// FullyDelete removes a content from the database
func (d *Draft) FullyDelete() error {
	return d.FullyDeleteTx(db.Con())
}

// FullyDeleteTx removes a content from the database using a transaction
func (d *Draft) FullyDeleteTx(tx sqalx.Node) error {
	if d == nil {
		return errors.New("content not instanced")
	}

	if d.ID == "" {
		return errors.New("content has not been saved")
	}

	stmt := "DELETE FROM blog_article_contents WHERE id=$1"
	_, err := tx.Exec(stmt, d.ID)

	return err
}

// Delete soft delete a draft.
func (d *Draft) Delete() error {
	return d.DeleteTx(db.Con())
}

// DeleteTx soft delete a draft using a transaction
func (d *Draft) DeleteTx(tx sqalx.Node) error {
	return d.doDelete(tx)
}

// doDelete performs a soft delete operation on a draft using an optional transaction
func (d *Draft) doDelete(tx sqalx.Node) error {
	if d == nil {
		return apierror.NewServerError("draft is not instanced")
	}

	if d.ID == "" {
		return apierror.NewServerError("cannot delete a non-persisted draft")
	}

	d.DeletedAt = db.Now()
	d.IsDraft = nil

	stmt := "UPDATE blog_article_contents SET deleted_at=:deleted_at WHERE id=:id"
	_, err := tx.NamedExec(stmt, d)
	return err
}

// IsZero checks if the object is either nil or don't have an ID
func (d *Draft) IsZero() bool {
	return d.ToContent().IsZero()
}
