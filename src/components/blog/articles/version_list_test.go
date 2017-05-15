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

func TestListVersion(t *testing.T) {
	testCases := []struct {
		description string
		f           func(t *testing.T)
	}{
		{"access", listVersionAccess},
		{"wrong input", listVersionWrongInput},
		{"Valid", listVersionValid},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, tc.f)
	}
}

func listVersionWrongInput(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	_, adminSession := testdata.NewAdminAuth(t)
	adminAuth := httptests.NewRequestAuth(adminSession)

	t.Run("parallel wrapper", func(t *testing.T) {
		testCases := []struct {
			description string
			code        int
			params      *articles.ListVersionParams
		}{
			{
				"Invalid Article ID",
				http.StatusBadRequest,
				&articles.ListVersionParams{ArticleID: "invalid"},
			},
			{
				"Unexisting Article ID",
				http.StatusNotFound,
				&articles.ListVersionParams{ArticleID: "9ea9cfcd-c455-4d22-848b-f87b32ae19f6"},
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				rec := callListVersion(t, tc.params, adminAuth)
				assert.Equal(t, tc.code, rec.Code)
			})
		}
	})
}

func listVersionAccess(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	_, userSession := testdata.NewAuth(t)
	userAuth := httptests.NewRequestAuth(userSession)
	a, _, adminSession := articletestdata.NewArticle(t, nil)
	adminAuth := httptests.NewRequestAuth(adminSession)

	t.Run("parallel wrapper", func(t *testing.T) {
		testCases := []struct {
			description string
			code        int
			auth        *httptests.RequestAuth
		}{
			{"Anonymous", http.StatusUnauthorized, nil},
			{"Logged user", http.StatusForbidden, userAuth},
			{"Admin", http.StatusOK, adminAuth},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				params := &articles.ListVersionParams{ArticleID: a.ID}
				rec := callListVersion(t, params, tc.auth)
				assert.Equal(t, tc.code, rec.Code)
			})
		}
	})
}

func listVersionValid(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	// Create 1 article with 1 version
	oneVersion, _, adminSession := articletestdata.NewArticle(t, nil)
	oneVersionList := []*articles.Version{oneVersion.Version}

	// Create 1 articles with 5 versions
	fiveVersion, _, _ := articletestdata.NewArticle(t, nil)
	fiveVersionList := []*articles.Version{fiveVersion.Version}
	for i := 0; i < 4; i++ {
		createdAt := time.Now().Add(time.Duration(i) * time.Second)
		base := &articles.Version{CreatedAt: &db.Time{Time: createdAt}}
		v := articletestdata.NewVersion(t, fiveVersion, base)
		fiveVersionList = append(fiveVersionList, v)
	}

	// Create 1 articles with 20 versions
	TwentyVersion, _, _ := articletestdata.NewArticle(t, nil)
	TwentyVersionList := []*articles.Version{TwentyVersion.Version}
	for i := 0; i < 19; i++ {
		createdAt := time.Now().Add(time.Duration(i) * time.Second)
		base := &articles.Version{CreatedAt: &db.Time{Time: createdAt}}
		v := articletestdata.NewVersion(t, TwentyVersion, base)
		TwentyVersionList = append(TwentyVersionList, v)
	}

	adminAuth := httptests.NewRequestAuth(adminSession)

	t.Run("parallel wrapper", func(t *testing.T) {
		testCases := []struct {
			description string
			article     *articles.Article
			totalWanted int
			orderedList []*articles.Version
		}{
			{"1 version", oneVersion, 1, oneVersionList},
			{"5 versions", fiveVersion, 5, fiveVersionList},
			{"20 versions", TwentyVersion, 20, TwentyVersionList},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				params := &articles.ListVersionParams{ArticleID: tc.article.ID}
				rec := callListVersion(t, params, adminAuth)
				require.Equal(t, http.StatusOK, rec.Code)

				var pld *articles.VersionsPayload
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}

				require.Equal(t, tc.totalWanted, len(pld.Results))
				for i, v := range pld.Results {
					assert.Equal(t, tc.orderedList[i].ID, v.ID)
					assert.Equal(t, tc.orderedList[i].Title, v.Title)
					assert.Equal(t, tc.orderedList[i].Subtitle, v.Subtitle)
					assert.Equal(t, tc.orderedList[i].Description, v.Description)
					assert.Equal(t, tc.orderedList[i].Content, v.Content)
				}
			})
		}
	})
}

func callListVersion(t *testing.T, params *articles.ListVersionParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: articles.Endpoints[articles.EndpointVersionList],
		Auth:     auth,
		Params:   params,
	}

	return httptests.NewRequest(t, ri)
}
