package education

import (
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-types/datetime"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
)

// Education represents an education
//go:generate api-cli generate model Education -t about_education -e Get,GetAny,JoinSQL
type Education struct {
	ID        string             `db:"id"`
	CreatedAt *datetime.DateTime `db:"created_at"`
	UpdatedAt *datetime.DateTime `db:"updated_at"`
	DeletedAt *datetime.DateTime `db:"deleted_at"`

	OrganizationID string `db:"organization_id"`
	Degree         string `db:"degree"`
	GPA            string `db:"gpa"`
	Location       string `db:"location"`
	Description    string `db:"description"`
	StartYear      int    `db:"start_year"`
	EndYear        *int   `db:"end_year"`

	// Embedded models
	*organizations.Organization `db:"org"`
}

// ListEducation represents a list of education
type ListEducation []*Education

// GetByID finds and returns an active education by ID
// Deleted object are not returned
func GetByID(q db.Queryable, id string) (*Education, error) {
	e := &Education{}
	stmt := `
	SELECT edu.*, ` + organizations.JoinSQL("org") + `
	FROM about_education edu
	JOIN about_organizations org
	  ON org.id = edu.organization_id
	WHERE edu.id=$1
	  AND edu.deleted_at IS NULL
	  AND org.deleted_at IS NULL
	LIMIT 1`
	err := q.Get(e, stmt, id)
	return e, apierror.NewFromSQL(err)
}

// GetAnyByID finds and returns an education by ID
// Deleted and orphan objects are returned
func GetAnyByID(q db.Queryable, id string) (*Education, error) {
	e := &Education{}
	stmt := `
	SELECT edu.*, ` + organizations.JoinSQL("org") + `
	FROM about_education edu
	JOIN about_organizations org
	  ON org.id = edu.organization_id
	WHERE edu.id=$1
	LIMIT 1`
	err := q.Get(e, stmt, id)
	return e, apierror.NewFromSQL(err)
}
