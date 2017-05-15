package articles_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/primitives/models/lifecycle"
	"github.com/Nivl/go-rest-tools/primitives/ptrs"
	"github.com/Nivl/go-rest-tools/security/auth/testdata"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles/articletestdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	testCases := []struct {
		description string
		f           func(t *testing.T)
	}{
		{"access data", accessData},
		{"update version", updateVersion},
		{"Update publish status", updatePublish},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, tc.f)
	}
}

// updateVersion tests that the version updates correctly
func updateVersion(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	a1, _, s1 := articletestdata.NewArticle(t, nil)
	auth := httptests.NewRequestAuth(s1)

	newVersion := articletestdata.NewVersion(t, a1, &articles.Version{
		Title:       "New title",
		Subtitle:    "New Subtitle",
		Description: "New Description",
		Content:     "New Content",
	})

	t.Run("parallel", func(t *testing.T) {
		testCases := []struct {
			description string
			code        int
			article     *articles.Article
			newVersion  *articles.Version
		}{
			{"Valid new version", http.StatusOK, a1, newVersion},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()

				params := &articles.UpdateParams{
					ID:      tc.article.ID,
					Version: tc.newVersion.ID,
				}
				rec := callUpdate(t, params, auth)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusOK {
					var pld *articles.Payload
					if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
						t.Fatal(err)
					}

					assert.Equal(t, newVersion.Title, pld.Title)
					assert.Equal(t, newVersion.Subtitle, pld.Subtitle)
					assert.Equal(t, newVersion.Description, pld.Description)
					assert.Equal(t, newVersion.Content, pld.Content)
				}
			})
		}
	})
}

// accessData tests that some wrong high level data are failing
func accessData(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	_, randomUserSession := testdata.NewAuth(t)
	a1, _, s1 := articletestdata.NewArticle(t, nil)

	t.Run("parallel", func(t *testing.T) {
		testCases := []struct {
			description string
			code        int
			auth        *httptests.RequestAuth
			params      *articles.UpdateParams
			article     *articles.Article
		}{
			{
				"As anonymous",
				http.StatusUnauthorized,
				nil,
				&articles.UpdateParams{ID: a1.ID},
				a1,
			},
			{
				"As logged user",
				http.StatusForbidden,
				httptests.NewRequestAuth(randomUserSession),
				&articles.UpdateParams{ID: a1.ID},
				a1,
			},
			{
				"Un-existing article",
				http.StatusNotFound,
				httptests.NewRequestAuth(s1),
				&articles.UpdateParams{ID: "1228f726-54b0-4b28-ad8b-fa7d3c1c37b7"},
				nil,
			},
			{
				"Invalid article uuid",
				http.StatusBadRequest,
				httptests.NewRequestAuth(s1),
				&articles.UpdateParams{ID: "invalid data"},
				nil,
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				rec := callUpdate(t, tc.params, tc.auth)
				assert.Equal(t, tc.code, rec.Code)
			})
		}
	})
}

// updatePublish tests that the article can be (un)published
func updatePublish(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	aToUnPublish, u, s := articletestdata.NewArticle(t, nil)
	aToStayPublished, _, _ := articletestdata.NewArticle(t, &articles.Article{User: u})
	aToPublish, _, _ := articletestdata.NewArticle(t, &articles.Article{User: u, PublishedAt: nil})
	aToStayUnPublished, _, _ := articletestdata.NewArticle(t, &articles.Article{User: u, PublishedAt: nil})

	auth := httptests.NewRequestAuth(s)

	t.Run("parallel", func(t *testing.T) {
		testCases := []struct {
			description string
			code        int
			article     *articles.Article
			publish     *bool
		}{
			{"Publish an article", http.StatusOK, aToPublish, ptrs.NewBool(true)},
			{"UnPublish an article", http.StatusOK, aToUnPublish, ptrs.NewBool(false)},
			{"Keep article published", http.StatusOK, aToStayPublished, nil},
			{"Keep article unpublished", http.StatusOK, aToStayUnPublished, nil},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()

				params := &articles.UpdateParams{
					ID:      tc.article.ID,
					Publish: tc.publish,
				}
				rec := callUpdate(t, params, auth)
				require.Equal(t, tc.code, rec.Code)

				var pld *articles.Payload
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}

				if tc.publish != nil {
					if *tc.publish {
						assert.NotNil(t, pld.PublishedAt)
					} else {
						assert.Nil(t, pld.PublishedAt)
					}
				} else {
					assert.Equal(t, tc.article.PublishedAt, pld.PublishedAt)
				}
			})
		}
	})
}

func callUpdate(t *testing.T, params *articles.UpdateParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: articles.Endpoints[articles.EndpointUpdate],
		Auth:     auth,
		Params:   params,
	}

	return httptests.NewRequest(t, ri)
}
