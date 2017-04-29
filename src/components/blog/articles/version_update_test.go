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
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles/articletestdata"
	"github.com/stretchr/testify/assert"
)

func TestUpdateVersion(t *testing.T) {
	testCases := []struct {
		description string
		f           func(t *testing.T)
	}{
		{"access data", updateVersionAccess},
		{"update with valid data", updateVersionValid},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, tc.f)
	}
}

// TestHandlerUpdate tests the update handler
func updateVersionAccess(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	_, userSession := testdata.NewAuth(t)
	userAuth := httptests.NewRequestAuth(userSession)
	a, _, adminSession := articletestdata.NewArticle(t, nil)
	adminAuth := httptests.NewRequestAuth(adminSession)

	t.Run("parallel wrapper", func(t *testing.T) {
		tests := []struct {
			description string
			code        int
			auth        *httptests.RequestAuth
			params      *articles.UpdateVersionParams
			article     *articles.Article
		}{
			{
				"As anonymous",
				http.StatusUnauthorized,
				nil,
				&articles.UpdateVersionParams{ID: a.Version.ID, ArticleID: a.ID},
				a,
			},
			{
				"As logged user",
				http.StatusUnauthorized,
				userAuth,
				&articles.UpdateVersionParams{ID: a.Version.ID, ArticleID: a.ID},
				a,
			},
			{
				"Un-existing article",
				http.StatusNotFound,
				adminAuth,
				&articles.UpdateVersionParams{ArticleID: "1228f726-54b0-4b28-ad8b-fa7d3c1c37b7", ID: "1228f726-54b0-4b28-ad8b-fa7d3c1c37b7"},
				nil,
			},
			{
				"Invalid article uuid",
				http.StatusBadRequest,
				adminAuth,
				&articles.UpdateVersionParams{ArticleID: "invalid data", ID: "invalid data"},
				nil,
			},
		}

		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				rec := callUpdateVersion(t, tc.params, tc.auth)
				assert.Equal(t, tc.code, rec.Code)
			})
		}
	})
}

// updateVersionValid tests the update handler actually updates
func updateVersionValid(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	t.Run("parallel wrapper", func(t *testing.T) {
		tests := []struct {
			description string
			code        int
			params      *articles.UpdateVersionParams
		}{
			{
				"update title",
				http.StatusOK,
				&articles.UpdateVersionParams{Title: "New Title"},
			},
			{
				"update subtitle",
				http.StatusOK,
				&articles.UpdateVersionParams{Subtitle: "New Subtitle"},
			},
			{
				"update description",
				http.StatusOK,
				&articles.UpdateVersionParams{Description: "New Description"},
			},
			{
				"update content",
				http.StatusOK,
				&articles.UpdateVersionParams{Content: "New Content"},
			},
			{
				"update all",
				http.StatusOK,
				&articles.UpdateVersionParams{Title: "New Title", Subtitle: "New Subtitle", Content: "New Content", Description: "New Description"},
			},
		}

		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()

				a, _, adminSession := articletestdata.NewArticle(t, nil)
				adminAuth := httptests.NewRequestAuth(adminSession)

				tc.params.ArticleID = a.ID
				tc.params.ID = a.Version.ID

				rec := callUpdateVersion(t, tc.params, adminAuth)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusOK {
					var pld *articles.VersionPayload
					if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
						t.Fatal(err)
					}

					if tc.params.Title != "" {
						assert.Equal(t, tc.params.Title, pld.Title)
					} else {
						assert.Equal(t, a.Title, pld.Title)
					}

					if tc.params.Subtitle != "" {
						assert.Equal(t, tc.params.Subtitle, pld.Subtitle)
					} else {
						assert.Equal(t, a.Subtitle, pld.Subtitle)
					}

					if tc.params.Description != "" {
						assert.Equal(t, tc.params.Description, pld.Description)
					} else {
						assert.Equal(t, a.Description, pld.Description)
					}

					if tc.params.Content != "" {
						assert.Equal(t, tc.params.Content, pld.Content)
					} else {
						assert.Equal(t, a.Content, pld.Content)
					}
				}
			})
		}
	})
}

func callUpdateVersion(t *testing.T, params *articles.UpdateVersionParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: articles.Endpoints[articles.EndpointVersionUpdate],
		Auth:     auth,
		Params:   params,
	}

	return httptests.NewRequest(t, ri)
}
