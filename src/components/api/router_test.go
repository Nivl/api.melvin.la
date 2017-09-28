package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/Nivl/go-rest-tools/logger/mocklogger"
	"github.com/Nivl/go-rest-tools/notifiers/mailer/mockmailer"
	"github.com/Nivl/go-rest-tools/notifiers/reporter/mockreporter"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/storage/filestorage/mockfilestorage"

	"github.com/Nivl/go-rest-tools/dependencies/mockdependencies"
	"github.com/melvin-laplanche/ml-api/src/components/api"
	"github.com/stretchr/testify/assert"
)

// Test that an un-existing route returns JSON and a 404
func TestRouteNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/404", nil)
	if err != nil {
		t.Fatal(err)
	}

	apiDeps := &mockdependencies.Dependencies{}
	apiDeps.On("FileStorage", mock.Anything).Return(&mockfilestorage.FileStorage{}, nil)
	apiDeps.On("DB").Return(&mockdb.Connection{})
	apiDeps.On("Mailer").Return(&mockmailer.Mailer{})
	apiDeps.On("Reporter").Return(&mockreporter.Reporter{})

	logger := &mocklogger.Logger{}
	logger.On("Errorf", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	apiDeps.On("Logger").Return(logger)

	rec := httptest.NewRecorder()
	api.GetRouter(apiDeps).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
}
