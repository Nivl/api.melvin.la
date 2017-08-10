package testexperience

import (
	"github.com/Nivl/go-rest-tools/storage/db"
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
		CreatedAt:      db.Now(),
		UpdatedAt:      db.Now(),
		JobTitle:       uniuri.New(),
		Description:    uniuri.New(),
		Location:       uniuri.New(),
		StartDate:      db.Today(),
		OrganizationID: org.ID,
		Organization:   org,
	}
}
