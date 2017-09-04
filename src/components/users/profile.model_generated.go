package users

// Code generated; DO NOT EDIT.

import (
	"errors"
	"fmt"

	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-rest-tools/storage/db"
	uuid "github.com/satori/go.uuid"
)

// JoinProfileSQL returns a string ready to be embed in a JOIN query
func JoinProfileSQL(prefix string) string {
	fields := []string{ "id", "created_at", "updated_at", "deleted_at", "user_id", "picture", "phone_number", "public_email", "linkedin_custom_url", "facebook_username", "twitter_username", "is_featured" }
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






// Save creates or updates the article depending on the value of the id using
// a transaction
func (p *Profile) Save(q db.Queryable) error {
	if p.ID == "" {
		return p.Create(q)
	}

	return p.Update(q)
}

// Create persists a profile in the database
func (p *Profile) Create(q db.Queryable) error {
	if p.ID != "" {
		return errors.New("cannot persist a profile that already has an ID")
	}

	return p.doCreate(q)
}

// doCreate persists a profile in the database using a Node
func (p *Profile) doCreate(q db.Queryable) error {
	p.ID = uuid.NewV4().String()
	p.UpdatedAt = db.Now()
	if p.CreatedAt == nil {
		p.CreatedAt = db.Now()
	}

	stmt := "INSERT INTO user_profiles (id, created_at, updated_at, deleted_at, user_id, picture, phone_number, public_email, linkedin_custom_url, facebook_username, twitter_username, is_featured) VALUES (:id, :created_at, :updated_at, :deleted_at, :user_id, :picture, :phone_number, :public_email, :linkedin_custom_url, :facebook_username, :twitter_username, :is_featured)"
	_, err := q.NamedExec(stmt, p)

  return apierror.NewFromSQL(err)
}

// Update updates most of the fields of a persisted profile
// Excluded fields are id, created_at, deleted_at, etc.
func (p *Profile) Update(q db.Queryable) error {
	if p.ID == "" {
		return errors.New("cannot update a non-persisted profile")
	}

	return p.doUpdate(q)
}

// doUpdate updates a profile in the database
func (p *Profile) doUpdate(q db.Queryable) error {
	if p.ID == "" {
		return errors.New("cannot update a non-persisted profile")
	}

	p.UpdatedAt = db.Now()

	stmt := "UPDATE user_profiles SET id=:id, created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, user_id=:user_id, picture=:picture, phone_number=:phone_number, public_email=:public_email, linkedin_custom_url=:linkedin_custom_url, facebook_username=:facebook_username, twitter_username=:twitter_username, is_featured=:is_featured WHERE id=:id"
	_, err := q.NamedExec(stmt, p)

	return apierror.NewFromSQL(err)
}

// Delete removes a profile from the database
func (p *Profile) Delete(q db.Queryable) error {
	if p.ID == "" {
		return errors.New("profile has not been saved")
	}

	stmt := "DELETE FROM user_profiles WHERE id=$1"
	_, err := q.Exec(stmt, p.ID)

	return err
}

// GetID returns the ID field
func (p *Profile) GetID() string {
	return p.ID
}

// SetID sets the ID field
func (p *Profile) SetID(id string) {
	p.ID = id
}

// IsZero checks if the object is either nil or don't have an ID
func (p *Profile) IsZero() bool {
	return p == nil || p.ID == ""
}