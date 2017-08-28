package users

// Code generated; DO NOT EDIT.

import (
	"testing"

		"github.com/Nivl/go-rest-tools/storage/db"

		"github.com/stretchr/testify/assert"

		"github.com/satori/go.uuid"

		"github.com/Nivl/go-rest-tools/storage/db/mockdb"
)


func TestProfileSaveNew(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*users.Profile")

	p := &Profile{}
	err := p.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, p.ID, "ID should have been set")
}

func TestProfileSaveExisting(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdate("*users.Profile")

	p := &Profile{}
	id := uuid.NewV4().String()
	p.ID = id
	err := p.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	mockDB.AssertExpectations(t)
	assert.Equal(t, id, p.ID, "ID should not have changed")
}

func TestProfileCreate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*users.Profile")

	p := &Profile{}
	err := p.Create(mockDB)

	assert.NoError(t, err, "Create() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, p.ID, "ID should have been set")
	assert.NotNil(t, p.CreatedAt, "CreatedAt should have been set")
	assert.NotNil(t, p.UpdatedAt, "UpdatedAt should have been set")
}

func TestProfileCreateWithID(t *testing.T) {
	mockDB := &mockdb.Queryable{}

	p := &Profile{}
	p.ID = uuid.NewV4().String()

	err := p.Create(mockDB)
	assert.Error(t, err, "Create() should have fail")
	mockDB.AssertExpectations(t)
}

func TestProfileCreateWithDate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*users.Profile")

	createdAt := db.Now().AddDate(0, 0, 1)
	p := &Profile{CreatedAt: createdAt}
	err := p.Create(mockDB)

	assert.NoError(t, err, "Create() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, p.ID, "ID should have been set")
	assert.True(t, p.CreatedAt.Equal(createdAt), "CreatedAt should not have been updated")
	assert.NotNil(t, p.UpdatedAt, "UpdatedAt should have been set")
}

func TestProfileCreateFail(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsertError("*users.Profile")

	p := &Profile{}
	err := p.Create(mockDB)

	assert.Error(t, err, "Create() should have fail")
	mockDB.AssertExpectations(t)
}


func TestProfileUpdate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdate("*users.Profile")

	p := &Profile{}
	p.ID = uuid.NewV4().String()
	err := p.Update(mockDB)

	assert.NoError(t, err, "Update() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, p.ID, "ID should have been set")
	assert.NotNil(t, p.UpdatedAt, "UpdatedAt should have been set")
}

func TestProfileUpdateWithoutID(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	p := &Profile{}
	err := p.Update(mockDB)

	assert.Error(t, err, "Update() should not have fail")
	mockDB.AssertExpectations(t)
}

func TestProfileUpdateFail(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdateError("*users.Profile")

	p := &Profile{}
	p.ID = uuid.NewV4().String()
	err := p.Update(mockDB)

	assert.Error(t, err, "Update() should have fail")
	mockDB.AssertExpectations(t)
}

func TestProfileDelete(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectDeletion()

	p := &Profile{}
	p.ID = uuid.NewV4().String()
	err := p.Delete(mockDB)

	assert.NoError(t, err, "Delete() should not have fail")
	mockDB.AssertExpectations(t)
}

func TestProfileDeleteWithoutID(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	p := &Profile{}
	err := p.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
	mockDB.AssertExpectations(t)
}

func TestProfileDeleteError(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectDeletionError()

	p := &Profile{}
	p.ID = uuid.NewV4().String()
	err := p.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
	mockDB.AssertExpectations(t)
}

func TestProfileGetID(t *testing.T) {
	p := &Profile{}
	p.ID = uuid.NewV4().String()
	assert.Equal(t, p.ID, p.GetID(), "GetID() did not return the right ID")
}

func TestProfileSetID(t *testing.T) {
	p := &Profile{}
	p.SetID(uuid.NewV4().String())
	assert.NotEmpty(t, p.ID, "SetID() did not set the ID")
}

func TestProfileIsZero(t *testing.T) {
	empty := &Profile{}
	assert.True(t, empty.IsZero(), "IsZero() should return true for empty struct")

	var nilStruct *Profile
	assert.True(t, nilStruct.IsZero(), "IsZero() should return true for nil struct")

	valid := &Profile{ID: uuid.NewV4().String()}
	assert.False(t, valid.IsZero(), "IsZero() should return false for valid struct")
}