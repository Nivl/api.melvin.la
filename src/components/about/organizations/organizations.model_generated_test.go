package organizations

// Code generated; DO NOT EDIT.

import (
	"testing"

		"github.com/Nivl/go-rest-tools/storage/db"

		"github.com/stretchr/testify/assert"

		"github.com/satori/go.uuid"

		"github.com/Nivl/go-rest-tools/storage/db/mockdb"
)


func TestOrganizationSaveNew(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*organizations.Organization")

	o := &Organization{}
	err := o.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, o.ID, "ID should have been set")
}

func TestOrganizationSaveExisting(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdate("*organizations.Organization")

	o := &Organization{}
	id := uuid.NewV4().String()
	o.ID = id
	err := o.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	mockDB.AssertExpectations(t)
	assert.Equal(t, id, o.ID, "ID should not have changed")
}

func TestOrganizationCreate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*organizations.Organization")

	o := &Organization{}
	err := o.Create(mockDB)

	assert.NoError(t, err, "Create() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, o.ID, "ID should have been set")
	assert.NotNil(t, o.CreatedAt, "CreatedAt should have been set")
	assert.NotNil(t, o.UpdatedAt, "UpdatedAt should have been set")
}

func TestOrganizationCreateWithID(t *testing.T) {
	mockDB := &mockdb.Queryable{}

	o := &Organization{}
	o.ID = uuid.NewV4().String()

	err := o.Create(mockDB)
	assert.Error(t, err, "Create() should have fail")
	mockDB.AssertExpectations(t)
}

func TestOrganizationCreateWithDate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*organizations.Organization")

	createdAt := db.Now().AddDate(0, 0, 1)
	o := &Organization{CreatedAt: createdAt}
	err := o.Create(mockDB)

	assert.NoError(t, err, "Create() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, o.ID, "ID should have been set")
	assert.True(t, o.CreatedAt.Equal(createdAt), "CreatedAt should not have been updated")
	assert.NotNil(t, o.UpdatedAt, "UpdatedAt should have been set")
}

func TestOrganizationCreateFail(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsertError("*organizations.Organization")

	o := &Organization{}
	err := o.Create(mockDB)

	assert.Error(t, err, "Create() should have fail")
	mockDB.AssertExpectations(t)
}


func TestOrganizationUpdate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdate("*organizations.Organization")

	o := &Organization{}
	o.ID = uuid.NewV4().String()
	err := o.Update(mockDB)

	assert.NoError(t, err, "Update() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, o.ID, "ID should have been set")
	assert.NotNil(t, o.UpdatedAt, "UpdatedAt should have been set")
}

func TestOrganizationUpdateWithoutID(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	o := &Organization{}
	err := o.Update(mockDB)

	assert.Error(t, err, "Update() should not have fail")
	mockDB.AssertExpectations(t)
}

func TestOrganizationUpdateFail(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdateError("*organizations.Organization")

	o := &Organization{}
	o.ID = uuid.NewV4().String()
	err := o.Update(mockDB)

	assert.Error(t, err, "Update() should have fail")
	mockDB.AssertExpectations(t)
}

func TestOrganizationDelete(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectDeletion()

	o := &Organization{}
	o.ID = uuid.NewV4().String()
	err := o.Delete(mockDB)

	assert.NoError(t, err, "Delete() should not have fail")
	mockDB.AssertExpectations(t)
}

func TestOrganizationDeleteWithoutID(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	o := &Organization{}
	err := o.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
	mockDB.AssertExpectations(t)
}

func TestOrganizationDeleteError(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectDeletionError()

	o := &Organization{}
	o.ID = uuid.NewV4().String()
	err := o.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
	mockDB.AssertExpectations(t)
}

func TestOrganizationGetID(t *testing.T) {
	o := &Organization{}
	o.ID = uuid.NewV4().String()
	assert.Equal(t, o.ID, o.GetID(), "GetID() did not return the right ID")
}

func TestOrganizationSetID(t *testing.T) {
	o := &Organization{}
	o.SetID(uuid.NewV4().String())
	assert.NotEmpty(t, o.ID, "SetID() did not set the ID")
}

func TestOrganizationIsZero(t *testing.T) {
	empty := &Organization{}
	assert.True(t, empty.IsZero(), "IsZero() should return true for empty struct")

	var nilStruct *Organization
	assert.True(t, nilStruct.IsZero(), "IsZero() should return true for nil struct")

	valid := &Organization{ID: uuid.NewV4().String()}
	assert.False(t, valid.IsZero(), "IsZero() should return false for valid struct")
}