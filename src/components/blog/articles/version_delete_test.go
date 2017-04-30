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

func TestDeleteVersion(t *testing.T) {
	testCases := []struct {
		description string
		f           func(t *testing.T)
	}{
		{"should be falling", deleteVersionFails},
		{"should succeed", deleteVersionSuccess},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, tc.f)
	}
}

func deleteVersionFails(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	_, userSession := testdata.NewAuth(t)
	userAuth := httptests.NewRequestAuth(userSession)

	a, _, adminSession := articletestdata.NewArticle(t, nil)
	v2 := articletestdata.NewVersion(t, a, nil)
	adminAuth := httptests.NewRequestAuth(adminSession)

	t.Run("parallel wrapper", func(t *testing.T) {
		tests := []struct {
			description string
			code        int
			params      *articles.DeleteVersionParams
			auth        *httptests.RequestAuth
		}{
			{
				"Anonymous",
				http.StatusUnauthorized,
				&articles.DeleteVersionParams{ID: v2.ID, ArticleID: a.ID},
				nil,
			},
			{
				"Logged user",
				http.StatusUnauthorized,
				&articles.DeleteVersionParams{ID: v2.ID, ArticleID: a.ID},
				userAuth,
			},
			{
				"Version in use by an article",
				http.StatusConflict,
				&articles.DeleteVersionParams{ID: a.Version.ID, ArticleID: a.ID},
				adminAuth,
			},
		}

		for _, tc := range tests {
			t.Run(tc.description, func(t *testing.T) {
				rec := callDeleteVersion(t, tc.params, tc.auth)
				assert.Equal(t, tc.code, rec.Code)
			})
		}
	})
}

func deleteVersionSuccess(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	a, _, adminSession := articletestdata.NewArticle(t, nil)
	v2 := articletestdata.NewVersion(t, a, nil)
	adminAuth := httptests.NewRequestAuth(adminSession)

	t.Run("parallel wrapper", func(t *testing.T) {
		tests := []struct {
			description string
			code        int
			params      *articles.DeleteVersionParams
			auth        *httptests.RequestAuth
		}{
			{
				"Removable version",
				http.StatusNoContent,
				&articles.DeleteVersionParams{ID: v2.ID, ArticleID: a.ID},
				adminAuth,
			},
		}

		for _, tc := range tests {
			t.Run(tc.description, func(t *testing.T) {
				rec := callDeleteVersion(t, tc.params, tc.auth)
				require.Equal(t, tc.code, rec.Code)

				v, err := articles.GetVersionByID(tc.params.ID)
				if err != nil {
					t.Fatal(err)
				}
				require.Nil(t, v)
			})
		}
	})
}

func callDeleteVersion(t *testing.T, params *articles.DeleteVersionParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: articles.Endpoints[articles.EndpointVersionDelete],
		Params:   params,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}
