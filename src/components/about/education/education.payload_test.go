package education_test

import (
	"testing"

	"github.com/melvin-laplanche/ml-api/src/components/about/education"

	"github.com/melvin-laplanche/ml-api/src/components/about/education/testeducation"
	"github.com/stretchr/testify/assert"
)

func TestExportPublic(t *testing.T) {
	exp := testeducation.New()
	pld := exp.ExportPublic()
	assert.Nil(t, pld.CreatedAt, "createdAt should not have been exported")
	assert.Nil(t, pld.UpdatedAt, "createdAt should not have been exported")
	assert.Nil(t, pld.DeletedAt, "deletedAt should not have been exported")

	assert.Equal(t, exp.ID, pld.ID, "ID should not have been changed")
	assert.Equal(t, exp.Degree, pld.Degree, "Degree should not have been changed")
	assert.Equal(t, exp.GPA, pld.GPA, "GPA should not have been changed")
	assert.Equal(t, exp.Description, pld.Description, "Description should not have been changed")
	assert.Equal(t, exp.Location, pld.Location, "Location should not have been changed")
}

func TestExportPublicNil(t *testing.T) {
	var exp *education.Education
	pld := exp.ExportPublic()
	assert.Nil(t, pld, "nil should export as nil")
}

func TestExportPrivate(t *testing.T) {
	exp := testeducation.New()
	pld := exp.ExportPrivate()
	assert.NotNil(t, pld.CreatedAt, "createdAt should have been exported")
	assert.NotNil(t, pld.UpdatedAt, "createdAt should have been exported")

	assert.Equal(t, exp.ID, pld.ID, "ID should not have been changed")
	assert.Equal(t, exp.Degree, pld.Degree, "Degree should not have been changed")
	assert.Equal(t, exp.GPA, pld.GPA, "GPA should not have been changed")
	assert.Equal(t, exp.Description, pld.Description, "Description should not have been changed")
	assert.Equal(t, exp.Location, pld.Location, "Location should not have been changed")
}

func TestExportPrivateNil(t *testing.T) {
	var exp *education.Education
	pld := exp.ExportPrivate()
	assert.Nil(t, pld, "nil should export as nil")
}

func TestExportListPublic(t *testing.T) {
	l := education.ListEducation{
		testeducation.New(),
		testeducation.New(),
		testeducation.New(),
		testeducation.New(),
		testeducation.New(),
	}
	pld := l.ExportPublic()

	if assert.Equal(t, len(l), len(pld.Results), "wrong number of exported results") {
		for i, r := range pld.Results {
			assert.Nil(t, r.CreatedAt, "createdAt should not have been exported")
			assert.Nil(t, r.UpdatedAt, "createdAt should not have been exported")
			assert.Nil(t, r.DeletedAt, "deletedAt should not have been exported")

			assert.Equal(t, l[i].ID, r.ID, "ID should not have been changed")
			assert.Equal(t, l[i].Degree, r.Degree, "Degree should not have been changed")
			assert.Equal(t, l[i].GPA, r.GPA, "GPA should not have been changed")
			assert.Equal(t, l[i].Description, r.Description, "Description should not have been changed")
			assert.Equal(t, l[i].Location, r.Location, "Location should not have been changed")
		}
	}
}

func TestExportListPrivate(t *testing.T) {
	l := education.ListEducation{
		testeducation.New(),
		testeducation.New(),
		testeducation.New(),
		testeducation.New(),
		testeducation.New(),
	}
	pld := l.ExportPrivate()

	if assert.Equal(t, len(l), len(pld.Results), "wrong number of exported results") {
		for i, r := range pld.Results {
			assert.NotNil(t, r.CreatedAt, "createdAt should have been exported")
			assert.NotNil(t, r.UpdatedAt, "createdAt should have been exported")

			assert.Equal(t, l[i].ID, r.ID, "ID should not have been changed")
			assert.Equal(t, l[i].Degree, r.Degree, "Degree should not have been changed")
			assert.Equal(t, l[i].GPA, r.GPA, "GPA should not have been changed")
			assert.Equal(t, l[i].Description, r.Description, "Description should not have been changed")
			assert.Equal(t, l[i].Location, r.Location, "Location should not have been changed")
		}
	}
}
