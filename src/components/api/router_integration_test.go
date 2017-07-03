// +build integration

package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/melvin-laplanche/ml-api/src/components/api"
	"github.com/stretchr/testify/assert"
)

// Test that an un-existing route returns JSON and a 404
func TestRouteNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/404", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	api.GetRouter().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
}
