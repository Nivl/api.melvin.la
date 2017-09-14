package testexperience

import (
	"testing"

	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-types/date"
	"github.com/Nivl/go-types/datetime"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/dchest/uniuri"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/satori/go.uuid"
)

// New returns a non persisted experience
func New() *experience.Experience {
	org := testorganizations.New()

	return &experience.Experience{
		ID:             uuid.NewV4().String(),
		CreatedAt:      datetime.Now(),
		UpdatedAt:      datetime.Now(),
		JobTitle:       uniuri.New(),
		Description:    uniuri.New(),
		Location:       uniuri.New(),
		StartDate:      date.Today(),
		OrganizationID: org.ID,
		Organization:   org,
	}
}

func NewPersisted(t *testing.T, dbCon db.Queryable, exp *experience.Experience) *experience.Experience {
	if exp == nil {
		exp = &experience.Experience{}
	}

	if exp.JobTitle == "" {
		exp.JobTitle = uniuri.New()
	}

	if exp.Description == "" {
		exp.Description = uniuri.New()
	}

	if exp.StartDate == nil {
		exp.StartDate = date.Today()
	}

	if exp.Organization != nil && exp.OrganizationID == "" {
		exp.OrganizationID = exp.Organization.ID
	}

	if exp.OrganizationID == "" {
		org := testorganizations.NewPersisted(t, dbCon, nil)
		exp.OrganizationID = org.ID
		exp.Organization = org
	}

	if err := exp.Create(dbCon); err != nil {
		t.Fatal(err)
	}

	lifecycle.SaveModels(t, exp)
	return exp
}
