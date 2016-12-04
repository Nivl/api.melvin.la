package articles_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/melvin-laplanche/ml-api/src/auth/authtest"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles/articlestest"
	"github.com/melvin-laplanche/ml-api/src/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestHandlerDelete(t *testing.T) {
	defer testhelpers.PurgeModels(t)

	// Random logged user
	randomUser, randomUserSession := authtest.NewAuth(t)

	// Published article with no draft
	a, u, s := articlestest.NewArticle(t, nil)

	tests := []struct {
		description string
		code        int
		params      *articles.HandlerGetParams
		auth        *testhelpers.RequestAuth
	}{
		{
			"Anonymous",
			http.StatusUnauthorized,
			&articles.HandlerGetParams{ID: a.ID},
			nil,
		},
		{
			"logged user",
			http.StatusNotFound,
			&articles.HandlerGetParams{ID: a.ID},
			testhelpers.NewRequestAuth(randomUserSession.ID, randomUser.ID),
		},
		{
			"Owner",
			http.StatusNoContent,
			&articles.HandlerGetParams{ID: a.ID},
			testhelpers.NewRequestAuth(s.ID, u.ID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerDelete(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusNoContent {
				exists, err := articles.Exists(tc.params.ID)
				if err != nil {
					t.Fatal(err)
				}
				assert.False(t, exists)
			}
		})
	}
}

func callHandlerDelete(t *testing.T, params *articles.HandlerGetParams, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: articles.Endpoints[articles.EndpointDelete],
		URI:      fmt.Sprintf("/blog/articles/%s", params.ID),
		Params:   params,
		Auth:     auth,
	}

	return testhelpers.NewRequest(ri)
}
