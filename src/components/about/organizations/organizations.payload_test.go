package organizations_test

import (
	"testing"

	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/stretchr/testify/assert"
)

func TestExportPublic(t *testing.T) {
	org := testorganizations.New()
	pld := org.ExportPublic()
	assert.Nil(t, pld.CreatedAt, "createdAt should not have been exported")
	assert.Nil(t, pld.UpdatedAt, "createdAt should not have been exported")
	assert.Nil(t, pld.DeletedAt, "deletedAt should not have been exported")

	assert.Equal(t, org.ID, pld.ID, "ID should not have been changed")
	assert.Equal(t, org.Name, pld.Name, "Name should not have been changed")
	assert.Equal(t, *org.ShortName, pld.ShortName, "ShortName should not have been changed")
	assert.Equal(t, *org.Logo, pld.Logo, "Logo should not have been changed")
	assert.Equal(t, *org.Website, pld.Website, "Website should not have been changed")
}

func TestExportPublicNil(t *testing.T) {
	var org *organizations.Organization
	pld := org.ExportPublic()
	assert.Nil(t, pld, "nil should export as nil")
}

func TestExportListPublic(t *testing.T) {
	l := organizations.Organizations{
		testorganizations.New(),
		testorganizations.New(),
		testorganizations.New(),
		testorganizations.New(),
		testorganizations.New(),
	}
	pld := l.ExportPublic()

	if assert.Equal(t, len(l), len(pld.Results), "wrong number of exported results") {
		for i, r := range pld.Results {
			assert.Nil(t, r.CreatedAt, "createdAt should not have been exported")
			assert.Nil(t, r.UpdatedAt, "createdAt should not have been exported")
			assert.Nil(t, r.DeletedAt, "deletedAt should not have been exported")

			assert.Equal(t, l[i].ID, r.ID, "ID should not have been changed")
			assert.Equal(t, l[i].Name, r.Name, "Name should not have been changed")
			assert.Equal(t, *l[i].ShortName, r.ShortName, "ShortName should not have been changed")
			assert.Equal(t, *l[i].Logo, r.Logo, "Logo should not have been changed")
			assert.Equal(t, *l[i].Website, r.Website, "Website should not have been changed")
		}
	}
}

func TestExportPrivate(t *testing.T) {
	org := testorganizations.New()
	pld := org.ExportPrivate()
	assert.NotNil(t, pld.CreatedAt, "createdAt should not have been exported")
	assert.NotNil(t, pld.UpdatedAt, "createdAt should not have been exported")

	assert.Equal(t, org.ID, pld.ID, "ID should not have been changed")
	assert.Equal(t, org.Name, pld.Name, "Name should not have been changed")
	assert.Equal(t, *org.ShortName, pld.ShortName, "ShortName should not have been changed")
	assert.Equal(t, *org.Logo, pld.Logo, "Logo should not have been changed")
	assert.Equal(t, *org.Website, pld.Website, "Website should not have been changed")
}

func TestExportPrivateNil(t *testing.T) {
	var org *organizations.Organization
	pld := org.ExportPrivate()
	assert.Nil(t, pld, "nil should export as nil")
}

func TestExportListPrivate(t *testing.T) {
	l := organizations.Organizations{
		testorganizations.New(),
		testorganizations.New(),
		testorganizations.New(),
		testorganizations.New(),
		testorganizations.New(),
	}
	pld := l.ExportPrivate()

	if assert.Equal(t, len(l), len(pld.Results), "wrong number of exported results") {
		for i, r := range pld.Results {
			assert.NotNil(t, r.CreatedAt, "createdAt should not have been exported")
			assert.NotNil(t, r.UpdatedAt, "createdAt should not have been exported")

			assert.Equal(t, l[i].ID, r.ID, "ID should not have been changed")
			assert.Equal(t, l[i].Name, r.Name, "Name should not have been changed")
			assert.Equal(t, *l[i].ShortName, r.ShortName, "ShortName should not have been changed")
			assert.Equal(t, *l[i].Logo, r.Logo, "Logo should not have been changed")
			assert.Equal(t, *l[i].Website, r.Website, "Website should not have been changed")
		}
	}
}
