package articles_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/api.melvin.la/api/app/testhelpers"
	"github.com/Nivl/api.melvin.la/api/components/blog/articles"
	"github.com/stretchr/testify/assert"
)

func TestHandlerAdd(t *testing.T) {
	tests := []struct {
		description string
		params      *articles.HandlerAddParams
		code        int
	}{
		{"No Title", &articles.HandlerAddParams{}, http.StatusBadRequest},
		{"As few params as possible", &articles.HandlerAddParams{Title: "My Super Article"}, http.StatusCreated},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerAdd(t, tc.params)
			assert.Equal(t, tc.code, rec.Code)
		})
	}
}

func callHandlerAdd(t *testing.T, params *articles.HandlerAddParams) *httptest.ResponseRecorder {
	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: articles.Endpoints[articles.EndpointAdd],
		URI:      "/blog/articles/",
		Params:   params,
	}

	return testhelpers.NewRequest(ri)
}
