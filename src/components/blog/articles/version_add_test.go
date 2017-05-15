package articles_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"time"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/primitives/models/lifecycle"
	"github.com/Nivl/go-rest-tools/security/auth/testdata"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles/articletestdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddVersion(t *testing.T) {
	testCases := []struct {
		description string
		f           func(t *testing.T)
	}{
		{"access", addVersionAccess},
		{"wrong input", addVersionWrongInput},
		{"Valid input", addVersionValid},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, tc.f)
	}
}

func addVersionValid(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	// Create a new article with a version in the future
	a, _, adminSession := articletestdata.NewArticle(t, nil)
	future := &db.Time{Time: time.Now().Add(1 * time.Hour)}
	a.Version.CreatedAt = future
	a.Version.Save()

	auth := httptests.NewRequestAuth(adminSession)
	params := &articles.AddVersionParams{ArticleID: a.ID}

	rec := callAddVersion(t, params, auth)
	require.Equal(t, http.StatusOK, rec.Code)

	var pld *articles.VersionPayload
	if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
		t.Fatal(err)
	}

	// Field that should be unique
	assert.NotEqual(t, a.Version.ID, pld.ID)
	assert.NotEqual(t, a.Version.CreatedAt.Unix(), pld.CreatedAt.Unix())

	// Field that should have been duplicated
	assert.Equal(t, a.Version.Title, pld.Title)
	assert.Equal(t, a.Version.Subtitle, pld.Subtitle)
	assert.Equal(t, a.Version.Content, pld.Content)
	assert.Equal(t, a.Version.Description, pld.Description)

	// We now make sure that the original version has not been altered
	v, err := articles.GetVersionByID(a.Version.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, a.Version.ID, v.ID)
	assert.Equal(t, a.Version.CreatedAt.Unix(), v.CreatedAt.Unix())
	assert.Equal(t, a.Version.UpdatedAt.Unix(), v.UpdatedAt.Unix())
	assert.Equal(t, a.Version.Title, v.Title)
	assert.Equal(t, a.Version.Subtitle, v.Subtitle)
	assert.Equal(t, a.Version.Content, v.Content)
	assert.Equal(t, a.Version.Description, v.Description)
}

func addVersionWrongInput(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	_, adminSession := testdata.NewAdminAuth(t)
	adminAuth := httptests.NewRequestAuth(adminSession)

	t.Run("parallel wrapper", func(t *testing.T) {
		testCases := []struct {
			description string
			code        int
			params      *articles.AddVersionParams
		}{
			{
				"Invalid article ID",
				http.StatusBadRequest,
				&articles.AddVersionParams{ArticleID: "invalid"},
			},
			{
				"un-existing article ID",
				http.StatusNotFound,
				&articles.AddVersionParams{ArticleID: "f069a2b5-0084-4f6d-b765-d13d616eb078"},
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				rec := callAddVersion(t, tc.params, adminAuth)
				assert.Equal(t, tc.code, rec.Code)
			})
		}
	})
}

func addVersionAccess(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	_, userSession := testdata.NewAuth(t)
	userAuth := httptests.NewRequestAuth(userSession)
	a, _, _ := articletestdata.NewArticle(t, nil)

	t.Run("parallel wrapper", func(t *testing.T) {
		testCases := []struct {
			description string
			code        int
			auth        *httptests.RequestAuth
		}{
			{"Anonymous", http.StatusUnauthorized, nil},
			{"Logged user", http.StatusForbidden, userAuth},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()

				params := &articles.AddVersionParams{ArticleID: a.ID}
				rec := callAddVersion(t, params, tc.auth)
				assert.Equal(t, tc.code, rec.Code)
			})
		}
	})
}

func callAddVersion(t *testing.T, params *articles.AddVersionParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: articles.Endpoints[articles.EndpointVersionAdd],
		Auth:     auth,
		Params:   params,
	}

	return httptests.NewRequest(t, ri)
}
