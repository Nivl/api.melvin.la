package experience

import (
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
)

// Experience represents a work experience
//go:generate api-cli generate model Experience -t about_experience -e Get
type Experience struct {
	ID        string   `db:"id"`
	CreatedAt *db.Time `db:"created_at"`
	UpdatedAt *db.Time `db:"updated_at"`
	DeletedAt *db.Time `db:"deleted_at"`

	OrganizationID string   `db:"organization_id"`
	JobTitle       string   `db:"job_title"`
	Location       string   `db:"location"`
	Description    string   `db:"description"`
	StartDate      *db.Time `db:"start_date"`
	EndDate        *db.Time `db:"end_date"`

	// Embedded models
	*organizations.Organization `db:"org"`
}

// List represents a list of experience
type List []*Experience

// GetByID finds and returns an active experience by ID
// Deleted object are not returned
func GetByID(q db.DB, id string) (*Experience, error) {
	e := &Experience{}
	stmt := `
	SELECT exp.*, ` + organizations.JoinSQL("org") + `
	FROM about_experience exp
	JOIN about_organizations org
	  ON org.id = exp.organization_id
	WHERE id=$1
	  AND exp.deleted_at IS NULL
	  AND org.deleted_at IS NULL
	LIMIT 1`
	err := q.Get(e, stmt, id)
	return e, apierror.NewFromSQL(err)
}
