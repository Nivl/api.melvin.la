package organizations

import (
	"github.com/Nivl/go-types/datetime"
	"github.com/Nivl/go-types/ptrs"
)

// Payload represents an article to be returned to the clients
type Payload struct {
	ID        string             `json:"id"`
	CreatedAt *datetime.DateTime `json:"created_at,omitempty"`
	UpdatedAt *datetime.DateTime `json:"updated_at,omitempty"`
	DeletedAt *datetime.DateTime `json:"deleted_at,omitempty"`
	Name      string             `json:"name"`
	ShortName string             `json:"short_name,omitempty"`
	Logo      string             `json:"logo,omitempty"`
	Website   string             `json:"website,omitempty"`
}

// ExportPublic returns an Organization payload trimmed from all sensitive
// information that can be
func (o *Organization) ExportPublic() *Payload {
	// It's OK to export a nil organization
	if o == nil {
		return nil
	}

	return &Payload{
		ID:        o.ID,
		Name:      o.Name,
		ShortName: ptrs.UnwrapString(o.ShortName),
		Logo:      ptrs.UnwrapString(o.Logo),
		Website:   ptrs.UnwrapString(o.Website),
	}
}

// ExportPrivate returns an Organization payload that can be
// safely returned to the clients
func (o *Organization) ExportPrivate() *Payload {
	// It's OK to export a nil organization
	if o == nil {
		return nil
	}

	pld := o.ExportPublic()
	pld.CreatedAt = o.CreatedAt
	pld.UpdatedAt = o.UpdatedAt
	pld.DeletedAt = o.DeletedAt
	return pld
}

// ListPayload represents a list of Organization that can be
// safely returned to the clients
type ListPayload struct {
	Results []*Payload `json:"results"`
}

// ExportPrivate returns a list of Organization as a payload that can be
// returned to the clients
func (o Organizations) ExportPrivate() *ListPayload {
	pld := &ListPayload{}
	pld.Results = make([]*Payload, len(o))
	for i, org := range o {
		pld.Results[i] = org.ExportPrivate()
	}
	return pld
}

// ExportPublic returns a list of Organization as a payload without any
// sensitive information
func (o Organizations) ExportPublic() *ListPayload {
	pld := &ListPayload{}
	pld.Results = make([]*Payload, len(o))
	for i, org := range o {
		pld.Results[i] = org.ExportPublic()
	}
	return pld
}
