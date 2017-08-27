package testrouter

import (
	"net/http"
	"testing"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ConflictTestParams represents the params needed by ConflictTest
type ConflictTestParams struct {
	Handler           router.RouteHandler
	HandlerParams     interface{}
	StructConflicting string
	FieldConflicting  string
}

// ConflictInsertTest test an handler and expects a 409 on an insert statement
func ConflictInsertTest(t *testing.T, p *ConflictTestParams) {
	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectInsertConflict(p.StructConflicting, p.FieldConflicting)

	// Mock the request & add expectations
	req := &mockrouter.HTTPRequest{}
	req.On("Params").Return(p.HandlerParams)

	// call the handler
	err := p.Handler(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusConflict, apiError.HTTPStatus())
	assert.Equal(t, p.FieldConflicting, apiError.Field())
}

// ConflictUpdateTest test am handler and expects a 409 on an update statement
func ConflictUpdateTest(t *testing.T, p *ConflictTestParams) {
	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGet(p.StructConflicting, func(args mock.Arguments) {
		obj := args.Get(0).(db.Model)
		obj.SetID("e9b51718-383b-42be-8720-0c8d7a99e978")
	})
	mockDB.ExpectUpdateConflict(p.StructConflicting, p.FieldConflicting)

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(p.HandlerParams)

	// call the handler
	err := p.Handler(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusConflict, apiError.HTTPStatus())
	assert.Equal(t, p.FieldConflicting, apiError.Field())
}
