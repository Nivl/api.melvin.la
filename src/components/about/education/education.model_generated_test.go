package education

// Code generated; DO NOT EDIT.

import (
	"testing"

		"github.com/Nivl/go-rest-tools/storage/db"

		"github.com/stretchr/testify/assert"

		"github.com/satori/go.uuid"

		"github.com/Nivl/go-rest-tools/storage/db/mockdb"
)


func TestEducationSaveNew(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*education.Education")

	e := &Education{}
	err := e.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, e.ID, "ID should have been set")
}

func TestEducationSaveExisting(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdate("*education.Education")

	e := &Education{}
	id := uuid.NewV4().String()
	e.ID = id
	err := e.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	mockDB.AssertExpectations(t)
	assert.Equal(t, id, e.ID, "ID should not have changed")
}

func TestEducationCreate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*education.Education")

	e := &Education{}
	err := e.Create(mockDB)

	assert.NoError(t, err, "Create() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, e.ID, "ID should have been set")
	assert.NotNil(t, e.CreatedAt, "CreatedAt should have been set")
	assert.NotNil(t, e.UpdatedAt, "UpdatedAt should have been set")
}

func TestEducationCreateWithID(t *testing.T) {
	mockDB := &mockdb.Queryable{}

	e := &Education{}
	e.ID = uuid.NewV4().String()

	err := e.Create(mockDB)
	assert.Error(t, err, "Create() should have fail")
	mockDB.AssertExpectations(t)
}

func TestEducationDoCreate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*education.Education")

	e := &Education{}
	err := e.doCreate(mockDB)

	assert.NoError(t, err, "doCreate() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, e.ID, "ID should have been set")
	assert.NotNil(t, e.CreatedAt, "CreatedAt should have been set")
	assert.NotNil(t, e.UpdatedAt, "UpdatedAt should have been set")
}

func TestEducationDoCreateWithDate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*education.Education")

	createdAt := db.Now().AddDate(0, 0, 1)
	e := &Education{CreatedAt: createdAt}
	err := e.doCreate(mockDB)

	assert.NoError(t, err, "doCreate() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, e.ID, "ID should have been set")
	assert.True(t, e.CreatedAt.Equal(createdAt), "CreatedAt should not have been updated")
	assert.NotNil(t, e.UpdatedAt, "UpdatedAt should have been set")
}

func TestEducationDoCreateFail(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsertError("*education.Education")

	e := &Education{}
	err := e.doCreate(mockDB)

	assert.Error(t, err, "doCreate() should have fail")
	mockDB.AssertExpectations(t)
}


func TestEducationUpdate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdate("*education.Education")

	e := &Education{}
	e.ID = uuid.NewV4().String()
	err := e.Update(mockDB)

	assert.NoError(t, err, "Update() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, e.ID, "ID should have been set")
	assert.NotNil(t, e.UpdatedAt, "UpdatedAt should have been set")
}

func TestEducationUpdateWithoutID(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	e := &Education{}
	err := e.Update(mockDB)

	assert.Error(t, err, "Update() should not have fail")
	mockDB.AssertExpectations(t)
}


func TestEducationDoUpdate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdate("*education.Education")

	e := &Education{}
	e.ID = uuid.NewV4().String()
	err := e.doUpdate(mockDB)

	assert.NoError(t, err, "doUpdate() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, e.ID, "ID should have been set")
	assert.NotNil(t, e.UpdatedAt, "UpdatedAt should have been set")
}

func TestEducationDoUpdateWithoutID(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	e := &Education{}
	err := e.doUpdate(mockDB)

	assert.Error(t, err, "doUpdate() should not have fail")
	mockDB.AssertExpectations(t)
}

func TestEducationDoUpdateFail(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdateError("*education.Education")

	e := &Education{}
	e.ID = uuid.NewV4().String()
	err := e.doUpdate(mockDB)

	assert.Error(t, err, "doUpdate() should have fail")
	mockDB.AssertExpectations(t)
}

func TestEducationDelete(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectDeletion()

	e := &Education{}
	e.ID = uuid.NewV4().String()
	err := e.Delete(mockDB)

	assert.NoError(t, err, "Delete() should not have fail")
	mockDB.AssertExpectations(t)
}

func TestEducationDeleteWithoutID(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	e := &Education{}
	err := e.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
	mockDB.AssertExpectations(t)
}

func TestEducationDeleteError(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectDeletionError()

	e := &Education{}
	e.ID = uuid.NewV4().String()
	err := e.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
	mockDB.AssertExpectations(t)
}

func TestEducationGetID(t *testing.T) {
	e := &Education{}
	e.ID = uuid.NewV4().String()
	assert.Equal(t, e.ID, e.GetID(), "GetID() did not return the right ID")
}

func TestEducationSetID(t *testing.T) {
	e := &Education{}
	e.SetID(uuid.NewV4().String())
	assert.NotEmpty(t, e.ID, "SetID() did not set the ID")
}

func TestEducationIsZero(t *testing.T) {
	empty := &Education{}
	assert.True(t, empty.IsZero(), "IsZero() should return true for empty struct")

	var nilStruct *Education
	assert.True(t, nilStruct.IsZero(), "IsZero() should return true for nil struct")

	valid := &Education{ID: uuid.NewV4().String()}
	assert.False(t, valid.IsZero(), "IsZero() should return false for valid struct")
}