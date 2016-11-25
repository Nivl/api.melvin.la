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
	"github.com/stretchr/testify/require"
)

func TestHandlerAdd(t *testing.T) {
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
			"Duplicate title (should generate a different slug)",
			http.StatusCreated,
			&articles.HandlerAddParams{Title: "My Super Article"},
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
		},
		{
			"All fields set",
			http.StatusCreated,
			&articles.HandlerAddParams{Title: "Title", Subtitle: "Subtitle", Description: "Description", Content: "Content"},
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerAdd(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusCreated {
				var a articles.PrivatePayload
				if err := json.NewDecoder(rec.Body).Decode(&a); err != nil {
					t.Fatal(err)
				}

				// Validate the article
				assert.NotEmpty(t, a.ID)
				assert.NotEmpty(t, a.Slug)
				assert.Nil(t, a.Draft)
				assert.Nil(t, a.PublishedAt)

				// Validate the article's content
				require.NotNil(t, a.Content)
				assert.Equal(t, tc.params.Title, a.Content.Title)
				assert.Equal(t, tc.params.Subtitle, a.Content.Subtitle)
				assert.Equal(t, tc.params.Description, a.Content.Description)
				assert.Equal(t, tc.params.Content, a.Content.Content)
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
