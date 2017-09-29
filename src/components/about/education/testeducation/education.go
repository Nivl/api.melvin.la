package testeducation

import (
	"math/rand"
	"testing"

	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/Nivl/go-types/datetime"
	"github.com/Nivl/go-types/ptrs"
	"github.com/dchest/uniuri"
	"github.com/melvin-laplanche/ml-api/src/components/about/education"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/satori/go.uuid"
)

// New returns a non persisted education
func New() *education.Education {
	org := testorganizations.New()

	return &education.Education{
		ID:             uuid.NewV4().String(),
		CreatedAt:      datetime.Now(),
		UpdatedAt:      datetime.Now(),
		Degree:         uniuri.New(),
		GPA:            ptrs.NewString(uniuri.NewLen(4)),
		Location:       ptrs.NewString(uniuri.New()),
		Description:    ptrs.NewString(uniuri.New()),
		StartYear:      rand.Intn(100) + 1950,
		OrganizationID: org.ID,
		Organization:   org,
	}
}

// NewPersisted returns a persisted education
func NewPersisted(t *testing.T, dbCon db.Queryable, edu *education.Education) *education.Education {
	if edu == nil {
		edu = &education.Education{}
	}

	if edu.Degree == "" {
		edu.Degree = uniuri.New()
	}

	if edu.GPA == nil {
		edu.GPA = ptrs.NewString(uniuri.NewLen(4))
	}

	if edu.Description == nil {
		edu.Description = ptrs.NewString(uniuri.New())
	}

	if edu.Location == nil {
		edu.Location = ptrs.NewString(uniuri.New())
	}

	if edu.StartYear == 0 {
		edu.StartYear = rand.Intn(100) + 1950
	}

	if edu.Organization != nil && edu.OrganizationID == "" {
		edu.OrganizationID = edu.Organization.ID
	}

	if edu.OrganizationID == "" {
		org := testorganizations.NewPersisted(t, dbCon, nil)
		edu.OrganizationID = org.ID
		edu.Organization = org
	}

	if err := edu.Create(dbCon); err != nil {
		t.Fatal(err)
	}

	lifecycle.SaveModels(t, edu)
	return edu
}
