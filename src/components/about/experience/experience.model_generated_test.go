package experience

// Code generated; DO NOT EDIT.

import (
	"testing"

		"github.com/stretchr/testify/assert"

		"github.com/satori/go.uuid"

		"github.com/Nivl/go-rest-tools/storage/db/mockdb"

	"github.com/Nivl/go-types/datetime"
)


func TestExperienceSaveNew(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*experience.Experience")

	e := &Experience{}
	err := e.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, e.ID, "ID should have been set")
}

func TestExperienceSaveExisting(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdate("*experience.Experience")

	e := &Experience{}
	id := uuid.NewV4().String()
	e.ID = id
	err := e.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	mockDB.AssertExpectations(t)
	assert.Equal(t, id, e.ID, "ID should not have changed")
}

func TestExperienceCreate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*experience.Experience")

	e := &Experience{}
	err := e.Create(mockDB)

	assert.NoError(t, err, "Create() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, e.ID, "ID should have been set")
	assert.NotNil(t, e.CreatedAt, "CreatedAt should have been set")
	assert.NotNil(t, e.UpdatedAt, "UpdatedAt should have been set")
}

func TestExperienceCreateWithID(t *testing.T) {
	mockDB := &mockdb.Queryable{}

	e := &Experience{}
	e.ID = uuid.NewV4().String()

	err := e.Create(mockDB)
	assert.Error(t, err, "Create() should have fail")
	mockDB.AssertExpectations(t)
}

func TestExperienceDoCreate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*experience.Experience")

	e := &Experience{}
	err := e.doCreate(mockDB)

	assert.NoError(t, err, "doCreate() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, e.ID, "ID should have been set")
	assert.NotNil(t, e.CreatedAt, "CreatedAt should have been set")
	assert.NotNil(t, e.UpdatedAt, "UpdatedAt should have been set")
}

func TestExperienceDoCreateWithDate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*experience.Experience")

	createdAt := datetime.Now().AddDate(0, 0, 1)
	e := &Experience{CreatedAt: createdAt}
	err := e.doCreate(mockDB)

	assert.NoError(t, err, "doCreate() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, e.ID, "ID should have been set")
	assert.True(t, e.CreatedAt.Equal(createdAt), "CreatedAt should not have been updated")
	assert.NotNil(t, e.UpdatedAt, "UpdatedAt should have been set")
}

func TestExperienceDoCreateFail(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsertError("*experience.Experience")

	e := &Experience{}
	err := e.doCreate(mockDB)

	assert.Error(t, err, "doCreate() should have fail")
	mockDB.AssertExpectations(t)
}


func TestExperienceUpdate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdate("*experience.Experience")

	e := &Experience{}
	e.ID = uuid.NewV4().String()
	err := e.Update(mockDB)

	assert.NoError(t, err, "Update() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, e.ID, "ID should have been set")
	assert.NotNil(t, e.UpdatedAt, "UpdatedAt should have been set")
}

func TestExperienceUpdateWithoutID(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	e := &Experience{}
	err := e.Update(mockDB)

	assert.Error(t, err, "Update() should not have fail")
	mockDB.AssertExpectations(t)
}


func TestExperienceDoUpdate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdate("*experience.Experience")

	e := &Experience{}
	e.ID = uuid.NewV4().String()
	err := e.doUpdate(mockDB)

	assert.NoError(t, err, "doUpdate() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, e.ID, "ID should have been set")
	assert.NotNil(t, e.UpdatedAt, "UpdatedAt should have been set")
}

func TestExperienceDoUpdateWithoutID(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	e := &Experience{}
	err := e.doUpdate(mockDB)

	assert.Error(t, err, "doUpdate() should not have fail")
	mockDB.AssertExpectations(t)
}

func TestExperienceDoUpdateFail(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdateError("*experience.Experience")

	e := &Experience{}
	e.ID = uuid.NewV4().String()
	err := e.doUpdate(mockDB)

	assert.Error(t, err, "doUpdate() should have fail")
	mockDB.AssertExpectations(t)
}

func TestExperienceDelete(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectDeletion()

	e := &Experience{}
	e.ID = uuid.NewV4().String()
	err := e.Delete(mockDB)

	assert.NoError(t, err, "Delete() should not have fail")
	mockDB.AssertExpectations(t)
}

func TestExperienceDeleteWithoutID(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	e := &Experience{}
	err := e.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
	mockDB.AssertExpectations(t)
}

func TestExperienceDeleteError(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectDeletionError()

	e := &Experience{}
	e.ID = uuid.NewV4().String()
	err := e.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
	mockDB.AssertExpectations(t)
}

func TestExperienceGetID(t *testing.T) {
	e := &Experience{}
	e.ID = uuid.NewV4().String()
	assert.Equal(t, e.ID, e.GetID(), "GetID() did not return the right ID")
}

func TestExperienceSetID(t *testing.T) {
	e := &Experience{}
	e.SetID(uuid.NewV4().String())
	assert.NotEmpty(t, e.ID, "SetID() did not set the ID")
}

func TestExperienceIsZero(t *testing.T) {
	empty := &Experience{}
	assert.True(t, empty.IsZero(), "IsZero() should return true for empty struct")

	var nilStruct *Experience
	assert.True(t, nilStruct.IsZero(), "IsZero() should return true for nil struct")

	valid := &Experience{ID: uuid.NewV4().String()}
	assert.False(t, valid.IsZero(), "IsZero() should return false for valid struct")
}