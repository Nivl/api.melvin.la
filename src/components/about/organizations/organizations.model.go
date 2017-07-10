package organizations

import (
	"github.com/Nivl/go-rest-tools/storage/db"
)

// Organization is a structure representing an article that can be saved in the database
//go:generate api-cli generate model Organization -t about_organizations
type Organization struct {
	ID        string   `db:"id"         json:"id"`
	CreatedAt *db.Time `db:"created_at" json:"created_at"`
	UpdatedAt *db.Time `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *db.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	Name      string   `db:"name"       json:"name"`
	ShortName *string  `db:"short_name" json:"short_name,omitempty"`
	Logo      *string  `db:"logo"       json:"logo,omitempty"`
	Website   *string  `db:"website"    json:"website,omitempty"`
}

// Export returns an Organization payload that can be
// safely returned to the clients
func (o *Organization) Export() *Organization {
	return o
}

// Organizations represents a list of Organization
type Organizations []*Organization

// Export returns a list of Organization as a payload that can be
// safely returned to the clients
func (o Organizations) Export() *ListPayload {
	return &ListPayload{Results: o}
}

// ListPayload represents a list of Organization that can be
// safely returned to the clients
type ListPayload struct {
	Results Organizations `json:"results"`
}
