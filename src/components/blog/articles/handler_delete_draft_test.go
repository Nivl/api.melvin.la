package articles_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/melvin-laplanche/ml-api/src/auth/authtest"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles/articlestest"
	"github.com/melvin-laplanche/ml-api/src/db"
	"github.com/melvin-laplanche/ml-api/src/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerDeleteDraft(t *testing.T) {
	defer testhelpers.PurgeModels(t)

	// Random logged user
	randomUser, randomUserSession := authtest.NewAuth(t)

	// Published article with a draft
	a, u, s := articlestest.NewArticle(t,
		&articles.Article{PublishedAt: db.Now(),
			Content: &articles.Content{Title: "Title", Content: "Content"},
			Draft:   &articles.Draft{Title: "Title Draft", Content: "Content Draft"},
		})

	tests := []struct {
		description string
		code        int
		params      *articles.HandlerDeleteDraftParams
		auth        *testhelpers.RequestAuth
	}{
		{
			"Anonymous",
			http.StatusUnauthorized,
			&articles.HandlerDeleteDraftParams{ID: a.ID},
			nil,
		},
		{
			"logged user",
			http.StatusNotFound,
			&articles.HandlerDeleteDraftParams{ID: a.ID},
			testhelpers.NewRequestAuth(randomUserSession.ID, randomUser.ID),
		},
		{
			"Owner",
			http.StatusNoContent,
			&articles.HandlerDeleteDraftParams{ID: a.ID},
			testhelpers.NewRequestAuth(s.ID, u.ID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerDeleteDraft(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusNoContent {
				a, err := articles.Get(tc.params.ID)
				if err != nil {
					t.Fatal(err)
				}

				// Let's make sure the article still exists
				require.NotNil(t, a)

				// Let's fetch the draft
				if err := a.FetchDraft(); err != nil {
					t.Fatal(err)
				}
				assert.Nil(t, a.Draft)
			}
		})
	}
}

func callHandlerDeleteDraft(t *testing.T, params *articles.HandlerDeleteDraftParams, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: articles.Endpoints[articles.EndpointDeleteDraft],
		URI:      fmt.Sprintf("/blog/articles/%s/draft", params.ID),
		Params:   params,
		Auth:     auth,
	}

	return testhelpers.NewRequest(ri)
}
