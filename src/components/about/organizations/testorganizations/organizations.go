package testorganizations

import (
	"testing"

	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/types/datetime"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/Nivl/go-rest-tools/types/ptrs"
	"github.com/dchest/uniuri"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	uuid "github.com/satori/go.uuid"
)

// New returns a non persisted organization
func New() *organizations.Organization {
	return &organizations.Organization{
		ID:        uuid.NewV4().String(),
		CreatedAt: datetime.Now(),
		UpdatedAt: datetime.Now(),
		Name:      uniuri.New(),
		ShortName: ptrs.NewString(uniuri.New()),
		Website:   ptrs.NewString("http://" + uniuri.New() + ".com"),
		Logo:      ptrs.NewString("http://via.placeholder.com/60x60"),
	}
}

// NewPersisted returns a new persisted organization with random data
func NewPersisted(t *testing.T, dbCon db.Queryable, org *organizations.Organization) *organizations.Organization {
	if org == nil {
		org = &organizations.Organization{}
	}

	if org.Name == "" {
		org.Name = uniuri.New()
	}
	if org.ShortName == nil {
		org.ShortName = ptrs.NewString(uniuri.New())
	}
	if org.Website == nil {
		org.Website = ptrs.NewString("http://" + uniuri.New() + ".com")
	}
	if org.Logo == nil {
		org.Logo = ptrs.NewString("http://via.placeholder.com/60x60")
	}

	if err := org.Create(dbCon); err != nil {
		t.Fatal(err)
	}

	lifecycle.SaveModels(t, org)
	return org
}
