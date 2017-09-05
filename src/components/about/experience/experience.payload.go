package experience

import (
	"github.com/Nivl/go-rest-tools/types/date"
	"github.com/Nivl/go-rest-tools/types/datetime"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
)

// Payload represents
type Payload struct {
	ID        string             `json:"id"`
	CreatedAt *datetime.DateTime `json:"created_at,omitempty"`
	UpdatedAt *datetime.DateTime `json:"updated_at,omitempty"`
	DeletedAt *datetime.DateTime `json:"deleted_at,omitempty"`

	JobTitle    string     `json:"job_title"`
	Location    string     `json:"location,omitempty"`
	Description string     `json:"description,omitempty"`
	StartDate   *date.Date `json:"start_date"`
	EndDate     *date.Date `json:"end_date,omitempty"`

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

// ListPayload represents a list of Experience that can be
// safely returned to the clients
type ListPayload struct {
	Results []*Payload `json:"results"`
}

// ExportPrivate returns a list of Experience as a payload that can be
// returned to the clients
func (e ListExperience) ExportPrivate() *ListPayload {
	pld := &ListPayload{}
	pld.Results = make([]*Payload, len(e))
	for i, exp := range e {
		pld.Results[i] = exp.ExportPrivate()
	}
	return pld
}

// ExportPublic returns a list of Experience as a payload without any
// sensitive information
func (e ListExperience) ExportPublic() *ListPayload {
	pld := &ListPayload{}
	pld.Results = make([]*Payload, len(e))
	for i, exp := range e {
		pld.Results[i] = exp.ExportPublic()
	}
	return pld
}
