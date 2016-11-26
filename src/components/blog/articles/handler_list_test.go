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

// TestHandlerList tests the List handler
func TestHandlerList(t *testing.T) {
	// Create 10 published articles
	for i := 0; i < 10; i++ {
		a, u, s := articles.NewTestArticle(t, nil)
		testhelpers.SaveModels(t, a, u, s)
	}

	u1, s1 := auth.NewTestAuth(t)
	testhelpers.SaveModels(t, u1, s1)

	// Create 10 unpublished articles
	for i := 0; i < 10; i++ {
		toSave := &articles.Article{PublishedAt: nil, User: u1, UserID: u1.ID}
		a, _, _ := articles.NewTestArticle(t, toSave)
		testhelpers.SaveModels(t, a)
	}

	defer testhelpers.PurgeModels(t)

	tests := []struct {
		description     string
		code            int
		auth            *testhelpers.RequestAuth
		nbArticleWanted int
	}{
		{"As anonymous", http.StatusOK, nil, 10},
		{"As logged user", http.StatusOK, testhelpers.NewRequestAuth(s1.ID, u1.ID), 10},
	}

	for i, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerList(t, tc.auth)
			assert.Equal(t, tc.code, rec.Code, "Test %d", i)

			if rec.Code == http.StatusOK {
				var body *articles.Payloads
				if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}

				// Validate the number of articles
				assert.Equal(t, tc.nbArticleWanted, len(body.Results), "Test %d", i)

				// validate the content of the articles
				for _, art := range body.Results {
					// Validate the visibility of the article
					assert.NotNil(t, art.PublishedAt, "Test %d", i)

					// validate the author
					require.NotNil(t, art.User, "Test %d", i)

					// Validate the content
					require.NotNil(t, art.Content, "Test %d", i)
					require.NotEmpty(t, art.Content.Title, "Test %d", i)
					assert.Nil(t, art.Draft, "Test %d", i)

				}
			}
		})
	}
}

func callHandlerList(t *testing.T, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: articles.Endpoints[articles.EndpointList],
		URI:      "/blog/articles",
		Auth:     auth,
	}

	return testhelpers.NewRequest(ri)
}
