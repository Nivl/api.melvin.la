package experience

import (
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
)

// Payload represents
type Payload struct {
	ID        string   `json:"id"`
	CreatedAt *db.Time `json:"created_at,omitempty"`
	UpdatedAt *db.Time `json:"updated_at,omitempty"`
	DeletedAt *db.Time `json:"deleted_at,omitempty"`

	JobTitle    string   `json:"job_title"`
	Location    string   `json:"location,omitempty"`
	Description string   `json:"description,omitempty"`
	StartDate   *db.Time `json:"start_date"`
	EndDate     *db.Time `json:"end_date,omitempty"`

	Organization *organizations.Payload `json:"organization,omitempty"`
}

// ExportPublic returns a Payload containing only the fields that are safe to
// be seen by anyone
func (e *Experience) ExportPublic() *Payload {
	// It's OK to export a nil experience
	if e == nil {
		return nil
	}

	return &Payload{
		ID:           e.ID,
		JobTitle:     e.JobTitle,
		Location:     e.Location,
		Description:  e.Description,
		StartDate:    e.StartDate,
		EndDate:      e.EndDate,
		Organization: e.Organization.ExportPublic(),
	}
}

// ExportPrivate returns a Payload containing all the fields
func (e *Experience) ExportPrivate() *Payload {
	// It's OK to export a nil experience
	if e == nil {
		return nil
	}

	pld := e.ExportPublic()
	pld.CreatedAt = e.CreatedAt
	pld.UpdatedAt = e.UpdatedAt
	pld.DeletedAt = e.DeletedAt
	pld.Organization = e.Organization.ExportPrivate()
	return pld
}
