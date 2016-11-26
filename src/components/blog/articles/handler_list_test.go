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

// ListTest tests the List handler
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
		description       string
		code              int
		auth              *testhelpers.RequestAuth
		publishedWanted   int
		unpublishedWanted int
	}{
		{
			"No params as anonymous",
			http.StatusOK,
			nil,
			10,
			0,
		},
		{
			"No params as logged user",
			http.StatusOK,
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
			10,
			0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerList(t, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusOK {
				var body *articles.Payloads
				if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}

				// Check the number of articles
				var nbPub int
				var nbUnPub int
				for _, art := range body.Results {
					require.NotNil(t, art.User)
					require.NotNil(t, art.Content)

					if art.PublishedAt != nil {
						nbPub++
					} else {
						nbUnPub++

						// An unpublished article is only listable by its author
						require.NotNil(t, tc.auth)
						assert.Equal(t, tc.auth.UserID, art.User.ID)
					}
				}
				assert.Equal(t, tc.publishedWanted, nbPub)
				assert.Equal(t, tc.unpublishedWanted, nbUnPub)
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
