package articles_test

import (
	"encoding/json"
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
		{"Title filled with spaces", &articles.HandlerAddParams{Title: "       "}, http.StatusBadRequest},
		{"As few params as possible", &articles.HandlerAddParams{Title: "My Super Article"}, http.StatusCreated},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerAdd(t, tc.params)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusCreated {
				var a articles.Article
				if err := json.NewDecoder(rec.Body).Decode(&a); err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, a.ID)
				assert.NotEmpty(t, a.Slug)
				assert.Equal(t, tc.params.Title, a.Title)
				if err := a.FullyDelete(); err != nil {
					t.Fatal(err)
				}
			}
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
