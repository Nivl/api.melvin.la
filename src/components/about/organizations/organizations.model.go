package organizations

import "github.com/Nivl/go-types/datetime"

// Organization is a structure representing an organization that can be saved in the database
//go:generate api-cli generate model Organization -t about_organizations
type Organization struct {
	ID        string             `db:"id"`
	CreatedAt *datetime.DateTime `db:"created_at"`
	UpdatedAt *datetime.DateTime `db:"updated_at"`
	DeletedAt *datetime.DateTime `db:"deleted_at"`
	Name      string             `db:"name"`
	ShortName *string            `db:"short_name"`
	Logo      *string            `db:"logo"`
	Website   *string            `db:"website"`
}

// Organizations represents a list of Organization
type Organizations []*Organization
