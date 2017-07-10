package testrouter

import (
	"net/http"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/stretchr/testify/assert"
)

// ConflictTestParams represents the params needed by ConflictTest
type ConflictTestParams struct {
	Handler           router.RouteHandler
	HandlerParams     interface{}
	StructConflicting string
	FieldConflicting  string
}

// ConflictTest test a handler and expects a 409
func ConflictTest(t *testing.T, p *ConflictTestParams) {
	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectInsertConflict(p.StructConflicting, p.FieldConflicting)

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(p.HandlerParams)

	// call the handler
	err := p.Handler(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := httperr.Convert(err)
	assert.Equal(t, http.StatusConflict, httpErr.Code())
	assert.Equal(t, p.FieldConflicting, httpErr.Field())
}
