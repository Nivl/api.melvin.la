package articles_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/primitives/models/lifecycle"
	"github.com/Nivl/go-rest-tools/security/auth/testdata"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	globalT := t
	defer lifecycle.PurgeModels(t)

	u1, s1 := testdata.NewAuth(t)

	tests := []struct {
		description string
		code        int
		params      *articles.AddParams
		auth        *httptests.RequestAuth
	}{
		{
			"No logged",
			http.StatusBadRequest,
			&articles.AddParams{},
			nil,
		},
		{
			"No Title",
			http.StatusBadRequest,
			&articles.AddParams{},
			httptests.NewRequestAuth(s1.ID, u1.ID),
		},
		{
			"Title filled with spaces",
			http.StatusBadRequest,
			&articles.AddParams{Title: "       "},
			httptests.NewRequestAuth(s1.ID, u1.ID),
		},
		{
			"As few params as possible",
			http.StatusCreated,
			&articles.AddParams{Title: "My Super Article"},
			httptests.NewRequestAuth(s1.ID, u1.ID),
		},
		{
			"Duplicate title",
			http.StatusConflict,
			&articles.AddParams{Title: "My Super Article"},
			httptests.NewRequestAuth(s1.ID, u1.ID),
		},
		{
			"All fields set",
			http.StatusCreated,
			&articles.AddParams{Title: "Title", Subtitle: "Subtitle", Description: "Description", Content: "Content"},
			httptests.NewRequestAuth(s1.ID, u1.ID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callAdd(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusCreated {
				var a articles.Payload
				if err := json.NewDecoder(rec.Body).Decode(&a); err != nil {
					t.Fatal(err)
				}

				// Validate the article
				assert.NotEmpty(t, a.ID)
				assert.NotEmpty(t, a.Slug)
				assert.Nil(t, a.PublishedAt)

				// Validate the article's content
				assert.Equal(t, tc.params.Title, a.Title)
				assert.Equal(t, tc.params.Subtitle, a.Subtitle)
				assert.Equal(t, tc.params.Description, a.Description)
				assert.Equal(t, tc.params.Content, a.Content)

				lifecycle.SaveModels(globalT, &articles.Article{ID: a.ID})
			}
		})
	}
}

func callAdd(t *testing.T, params *articles.AddParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: articles.Endpoints[articles.EndpointAdd],
		Params:   params,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}
