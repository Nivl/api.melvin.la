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
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	testCases := []struct {
		description string
		f           func(t *testing.T)
	}{
		{"access data", getAccessData},
		{"returned data", getReturnedData},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, tc.f)
	}
}

func getAccessData(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	aPublished, _, _ := articletestdata.NewArticle(t, nil)
	aUnPublished, _, adminSession := articletestdata.NewArticle(t, &articles.Article{PublishedAt: nil})
	adminAuth := httptests.NewRequestAuth(adminSession)

	_, userSession := testdata.NewAuth(t)
	userAuth := httptests.NewRequestAuth(userSession)

	t.Run("parallel", func(t *testing.T) {
		testCases := []struct {
			description string
			code        int
			auth        *httptests.RequestAuth
			ArticleID   string
		}{
			{"anonymous accessing unpublished", http.StatusNotFound, nil, aUnPublished.ID},
			{"user accessing unpublished", http.StatusNotFound, userAuth, aUnPublished.ID},
			{"admin accessing unpublished", http.StatusOK, adminAuth, aUnPublished.ID},

			{"Anonymous accessing published", http.StatusOK, nil, aPublished.ID},
			{"user accessing published", http.StatusOK, userAuth, aPublished.ID},
			{"admin accessing published", http.StatusOK, adminAuth, aPublished.ID},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()

				params := &articles.GetParams{ID: tc.ArticleID}
				rec := callGet(t, params, tc.auth)
				require.Equal(t, tc.code, rec.Code)
			})
		}
	})
}

func getReturnedData(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	aPublished, _, adminSession := articletestdata.NewArticle(t, nil)
	adminAuth := httptests.NewRequestAuth(adminSession)

	_, userSession := testdata.NewAuth(t)
	userAuth := httptests.NewRequestAuth(userSession)

	// Everyone (admin, user, anon) gets the same data
	t.Run("parallel", func(t *testing.T) {
		testCases := []struct {
			description string
			auth        *httptests.RequestAuth
			article     *articles.Article
		}{
			{"anonymous", nil, aPublished},
			{"user", userAuth, aPublished},
			{"admin", adminAuth, aPublished},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()

				params := &articles.GetParams{ID: tc.article.ID}
				rec := callGet(t, params, tc.auth)
				require.Equal(t, http.StatusOK, rec.Code)

				if rec.Code == http.StatusOK {
					var pld articles.Payload
					if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
						t.Fatal(err)
					}

					assert.Equal(t, tc.article.ID, pld.ID)
					assert.Equal(t, tc.article.Slug, pld.Slug)

					require.NotNil(t, pld.User)
					assert.Equal(t, tc.article.UserID, pld.User.ID)

					assert.Equal(t, tc.article.Version.Title, pld.Title)
					assert.Equal(t, tc.article.Version.Subtitle, pld.Subtitle)
					assert.Equal(t, tc.article.Version.Description, pld.Description)
					assert.Equal(t, tc.article.Version.Content, pld.Content)
				}
			})
		}
	})
}

func callGet(t *testing.T, params *articles.GetParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: articles.Endpoints[articles.EndpointGet],
		Params:   params,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}
