package education

import (
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/types/ptrs"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
)

// Payload represents
type Payload struct {
	ID        string   `json:"id"`
	CreatedAt *db.Time `json:"created_at,omitempty"`
	UpdatedAt *db.Time `json:"updated_at,omitempty"`
	DeletedAt *db.Time `json:"deleted_at,omitempty"`

	Degree      string `json:"degree"`
	GPA         string `json:"gpa,omitempty"`
	Location    string `json:"location,omitempty"`
	Description string `json:"description,omitempty"`
	StartYear   int    `json:"start_year"`
	EndYear     int    `json:"end_year,omitempty"`

	Organization *organizations.Payload `json:"organization,omitempty"`
}

// ExportPublic returns a Payload containing only the fields that are safe to
// be seen by anyone
func (e *Education) ExportPublic() *Payload {
	// It's OK to export a nil education
	if e == nil {
		return nil
	}

	return &Payload{
		ID:           e.ID,
		Degree:       e.Degree,
		GPA:          e.GPA,
		Location:     e.Location,
		Description:  e.Description,
		StartYear:    e.StartYear,
		EndYear:      ptrs.UnwrapInt(e.EndYear),
		Organization: e.Organization.ExportPublic(),
	}
}

// ExportPrivate returns a Payload containing all the fields
func (e *Education) ExportPrivate() *Payload {
	// It's OK to export a nil education
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

// ListPayload represents a list of Education that can be
// safely returned to the clients
type ListPayload struct {
	Results []*Payload `json:"results"`
}

// ExportPrivate returns a list of Education as a payload that can be
// returned to the clients
func (e ListEducation) ExportPrivate() *ListPayload {
	pld := &ListPayload{}
	pld.Results = make([]*Payload, len(e))
	for i, exp := range e {
		pld.Results[i] = exp.ExportPrivate()
	}
	return pld
}

// ExportPublic returns a list of Education as a payload without any
// sensitive information
func (e ListEducation) ExportPublic() *ListPayload {
	pld := &ListPayload{}
	pld.Results = make([]*Payload, len(e))
	for i, exp := range e {
		pld.Results[i] = exp.ExportPublic()
	}
	return pld
}
