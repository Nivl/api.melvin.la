package organizations

import (
	"github.com/Nivl/go-rest-tools/primitives/ptrs"
	"github.com/Nivl/go-rest-tools/storage/db"
)

// Payload represents an article to be returned to the clients
type Payload struct {
	ID        string   `db:"id"         json:"id"`
	CreatedAt *db.Time `json:"created_at"`
	UpdatedAt *db.Time `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *db.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	Name      string   `db:"name"       json:"name"`
	ShortName string   `db:"short_name" json:"short_name,omitempty"`
	Logo      string   `db:"logo"       json:"logo,omitempty"`
	Website   string   `db:"website"    json:"website,omitempty"`
}

// ExportPublic returns an Organization payload trimmed from all sensitive
// information that can be
func (o *Organization) ExportPublic() *Payload {
	pld := &Payload{
		ID:        o.ID,
		Name:      o.Name,
		ShortName: ptrs.UnwrapString(o.ShortName),
		Logo:      ptrs.UnwrapString(o.Logo),
		Website:   ptrs.UnwrapString(o.Website),
	}
	return pld
}

// ExportPrivate returns an Organization payload that can be
// safely returned to the clients
func (o *Organization) ExportPrivate() *Payload {
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
