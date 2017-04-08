package articles

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
	fields := []string{ "id", "slug", "created_at", "updated_at", "deleted_at", "published_at", "user_id", "current_version" }
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

// GetByID finds and returns an active article by ID
func GetByID(id string) (*Article, error) {
	a := &Article{}
	stmt := "SELECT * from blog_articles WHERE id=$1 and deleted_at IS NULL LIMIT 1"
	err := db.Get(a, stmt, id)
	// We want to return nil if a article is not found
	if a.ID == "" {
		return nil, err
	}
	return a, err
}

// Exists checks if a article exists for a specific ID
func Exists(id string) (bool, error) {
	exists := false
	stmt := "SELECT exists(SELECT 1 FROM blog_articles WHERE id=$1 and deleted_at IS NULL)"
	err := db.Writer.Get(&exists, stmt, id)
	return exists, err
}

// Save creates or updates the article depending on the value of the id
func (a *Article) Save() error {
	return a.SaveQ(db.Writer)
}

// SaveQ creates or updates the article depending on the value of the id using
// a transaction
func (a *Article) SaveQ(q db.Queryable) error {
	if a == nil {
		return httperr.NewServerError("article is not instanced")
	}

	if a.ID == "" {
		return a.CreateQ(q)
	}

	return a.UpdateQ(q)
}

// Create persists a article in the database
func (a *Article) Create() error {
	return a.CreateQ(db.Writer)
}

// CreateQ persists a article in the database
func (a *Article) CreateQ(q db.Queryable) error {
	if a == nil {
		return httperr.NewServerError("article is not instanced")
	}

	if a.ID != "" {
		return httperr.NewServerError("cannot persist a article that already has an ID")
	}

	return a.doCreate(q)
}

// doCreate persists a article in the database using a Node
func (a *Article) doCreate(q db.Queryable) error {
	if a == nil {
		return errors.New("article not instanced")
	}

	a.ID = uuid.NewV4().String()
	a.UpdatedAt = db.Now()
	if a.CreatedAt == nil {
		a.CreatedAt = db.Now()
	}

	stmt := "INSERT INTO blog_articles (id, slug, created_at, updated_at, deleted_at, published_at, user_id, current_version) VALUES (:id, :slug, :created_at, :updated_at, :deleted_at, :published_at, :user_id, :current_version)"
	_, err := q.NamedExec(stmt, a)

  return err
}

// Update updates most of the fields of a persisted article.
// Excluded fields are id, created_at, deleted_at, etc.
func (a *Article) Update() error {
	return a.UpdateQ(db.Writer)
}

// UpdateQ updates most of the fields of a persisted article using a transaction
// Excluded fields are id, created_at, deleted_at, etc.
func (a *Article) UpdateQ(q db.Queryable) error {
	if a == nil {
		return httperr.NewServerError("article is not instanced")
	}

	if a.ID == "" {
		return httperr.NewServerError("cannot update a non-persisted article")
	}

	return a.doUpdate(q)
}

// doUpdate updates a article in the database using an optional transaction
func (a *Article) doUpdate(q db.Queryable) error {
	if a == nil {
		return httperr.NewServerError("article is not instanced")
	}

	if a.ID == "" {
		return httperr.NewServerError("cannot update a non-persisted article")
	}

	a.UpdatedAt = db.Now()

	stmt := "UPDATE blog_articles SET id=:id, slug=:slug, created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, published_at=:published_at, user_id=:user_id, current_version=:current_version WHERE id=:id"
	_, err := q.NamedExec(stmt, a)

	return err
}

// FullyDelete removes a article from the database
func (a *Article) FullyDelete() error {
	return a.FullyDeleteQ(db.Writer)
}

// FullyDeleteQ removes a article from the database using a transaction
func (a *Article) FullyDeleteQ(q db.Queryable) error {
	if a == nil {
		return errors.New("article not instanced")
	}

	if a.ID == "" {
		return errors.New("article has not been saved")
	}

	stmt := "DELETE FROM blog_articles WHERE id=$1"
	_, err := q.Exec(stmt, a.ID)

	return err
}

// Delete soft delete a article.
func (a *Article) Delete() error {
	return a.DeleteQ(db.Writer)
}

// DeleteQ soft delete a article using a transaction
func (a *Article) DeleteQ(q db.Queryable) error {
	return a.doDelete(q)
}

// doDelete performs a soft delete operation on a article using an optional transaction
func (a *Article) doDelete(q db.Queryable) error {
	if a == nil {
		return httperr.NewServerError("article is not instanced")
	}

	if a.ID == "" {
		return httperr.NewServerError("cannot delete a non-persisted article")
	}

	a.DeletedAt = db.Now()

	stmt := "UPDATE blog_articles SET deleted_at = $2 WHERE id=$1"
	_, err := q.Exec(stmt, a.ID, a.DeletedAt)
	return err
}

// IsZero checks if the object is either nil or don't have an ID
func (a *Article) IsZero() bool {
	return a == nil || a.ID == ""
}