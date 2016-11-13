package articles_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles"
	"github.com/melvin-laplanche/ml-api/src/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestHandlerAdd(t *testing.T) {
	globalT := t
	defer testhelpers.PurgeModels(t)

	u1, s1 := auth.NewTestAuth(t)
	testhelpers.SaveModels(t, u1, s1)

	tests := []struct {
		description string
		code        int
		params      *articles.HandlerAddParams
		auth        *testhelpers.RequestAuth
	}{
		{
			"No logged",
			http.StatusBadRequest,
			&articles.HandlerAddParams{},
			nil,
		},
		{
			"No Title",
			http.StatusBadRequest,
			&articles.HandlerAddParams{},
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
		},
		{
			"Title filled with spaces",
			http.StatusBadRequest,
			&articles.HandlerAddParams{Title: "       "},
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
		},
		{
			"As few params as possible",
			http.StatusCreated,
			&articles.HandlerAddParams{Title: "My Super Article"},
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
		},
		{
			"Duplicate title",
			http.StatusCreated,
			&articles.HandlerAddParams{Title: "My Super Article"},
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerAdd(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusCreated {
				var a articles.PublicPayload
				if err := json.NewDecoder(rec.Body).Decode(&a); err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, a.ID)
				assert.NotEmpty(t, a.Slug)
				assert.Equal(t, tc.params.Title, a.Title)
				testhelpers.SaveModels(globalT, &articles.Article{ID: a.ID})
			}
		})
	}
}

func callHandlerAdd(t *testing.T, params *articles.HandlerAddParams, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: articles.Endpoints[articles.EndpointAdd],
		URI:      "/blog/articles/",
		Params:   params,
		Auth:     auth,
	}

	return testhelpers.NewRequest(ri)
}
