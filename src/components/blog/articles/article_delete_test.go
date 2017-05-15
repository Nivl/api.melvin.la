package articles_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/primitives/models/lifecycle"
	"github.com/Nivl/go-rest-tools/security/auth/testdata"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles/articletestdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	defer lifecycle.PurgeModels(t)

	_, userSession := testdata.NewAuth(t)

	published, _, adminSession := articletestdata.NewArticle(t, nil)
	unPublished, _, _ := articletestdata.NewArticle(t, &articles.Article{PublishedAt: nil})

	tests := []struct {
		description string
		code        int
		params      *articles.DeleteParams
		auth        *httptests.RequestAuth
	}{
		{
			"Anonymous",
			http.StatusUnauthorized,
			&articles.DeleteParams{ID: published.ID},
			nil,
		},
		{
			"logged user",
			http.StatusForbidden,
			&articles.DeleteParams{ID: published.ID},
			httptests.NewRequestAuth(userSession),
		},
		{
			"Admin published",
			http.StatusConflict,
			&articles.DeleteParams{ID: published.ID},
			httptests.NewRequestAuth(adminSession),
		},
		{
			"Admin not published",
			http.StatusNoContent,
			&articles.DeleteParams{ID: unPublished.ID},
			httptests.NewRequestAuth(adminSession),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callDelete(t, tc.params, tc.auth)
			require.Equal(t, tc.code, rec.Code)

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

func callDelete(t *testing.T, params *articles.DeleteParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: articles.Endpoints[articles.EndpointDelete],
		Params:   params,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}
