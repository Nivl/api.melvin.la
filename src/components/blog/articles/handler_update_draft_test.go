package articles_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/melvin-laplanche/ml-api/src/auth/authtest"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles/articlestest"
	"github.com/melvin-laplanche/ml-api/src/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHandlerUpdate tests the update handler
func TestHandlerUpdateDraft(t *testing.T) {
	defer testhelpers.PurgeModels(t)

	randomUser, randomUserSession := authtest.NewAuth(t)

	a1, u1, s1 := articlestest.NewArticle(t, nil)
	a2, u2, s2 := articlestest.NewArticle(t, nil)

	tests := []struct {
		description string
		code        int
		auth        *testhelpers.RequestAuth
		params      *articles.HandlerUpdateDraftParams
		article     *articles.Article
	}{
		{
			"As anonymous",
			http.StatusUnauthorized,
			nil,
			&articles.HandlerUpdateDraftParams{ID: a1.ID},
			a1,
		},
		{
			"As logged user",
			http.StatusNotFound,
			testhelpers.NewRequestAuth(randomUserSession.ID, randomUser.ID),
			&articles.HandlerUpdateDraftParams{ID: a1.ID},
			a1,
		},
		{
			"Un-existing article",
			http.StatusNotFound,
			testhelpers.NewRequestAuth(randomUserSession.ID, randomUser.ID),
			&articles.HandlerUpdateDraftParams{ID: "1228f726-54b0-4b28-ad8b-fa7d3c1c37b7"},
			nil,
		},
		{
			"Invalid article uuid",
			http.StatusBadRequest,
			testhelpers.NewRequestAuth(randomUserSession.ID, randomUser.ID),
			&articles.HandlerUpdateDraftParams{ID: "invalid data"},
			nil,
		},
		{
			"As author - changing draft",
			http.StatusOK,
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
			&articles.HandlerUpdateDraftParams{ID: a1.ID, Content: "New Content", Subtitle: "New Subtitle"},
			a1,
		},
		{
			"As author - Promote",
			http.StatusOK,
			testhelpers.NewRequestAuth(s2.ID, u2.ID),
			&articles.HandlerUpdateDraftParams{ID: a2.ID, Promote: true},
			a2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerUpdateDraft(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if testhelpers.Is2XX(rec.Code) {
				var pld *articles.Payload
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}

				toTest := pld.Draft
				if tc.params.Promote {
					require.Nil(t, pld.Draft)
					require.NotNil(t, pld.Content)
					toTest = pld.Content
				}

				if tc.params.Title != "" {
					assert.Equal(t, tc.params.Title, toTest.Title)
				} else {
					assert.Equal(t, tc.article.Content.Title, pld.Content.Title)
				}

				if tc.params.Subtitle != "" {
					assert.Equal(t, tc.params.Subtitle, toTest.Subtitle)
				} else {
					assert.Equal(t, tc.article.Content.Subtitle, pld.Content.Subtitle)
				}

				if tc.params.Description != "" {
					assert.Equal(t, tc.params.Description, toTest.Description)
				} else {
					assert.Equal(t, tc.article.Content.Description, pld.Content.Description)
				}

				if tc.params.Content != "" {
					assert.Equal(t, tc.params.Content, toTest.Content)
				} else {
					assert.Equal(t, tc.article.Content.Content, pld.Content.Content)
				}
			}
		})
	}
}

func callHandlerUpdateDraft(t *testing.T, params *articles.HandlerUpdateDraftParams, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: articles.Endpoints[articles.EndpointUpdateDraft],
		URI:      fmt.Sprintf("/blog/articles/%s/draft", params.ID),
		Auth:     auth,
		Params:   params,
	}

	return testhelpers.NewRequest(ri)
}
