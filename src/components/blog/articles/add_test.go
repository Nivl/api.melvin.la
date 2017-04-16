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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		description string
		f           func(t *testing.T)
	}{
		{"access data", addAccessData},
		{"invalid title", addInvalidTitle},
		{"valid data", addValid},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, tc.f)
	}
}

func addInvalidTitle(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	_, s := testdata.NewAuth(t)
	auth := httptests.NewRequestAuth(s)

	t.Run("parallel", func(t *testing.T) {
		testCases := []struct {
			description string
			code        int
			params      *articles.AddParams
		}{
			{
				"No Title",
				http.StatusBadRequest,
				&articles.AddParams{},
			},
			{
				"Title filled with spaces",
				http.StatusBadRequest,
				&articles.AddParams{Title: "       "},
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()

				rec := callAdd(t, tc.params, auth)
				assert.Equal(t, tc.code, rec.Code)
			})
		}
	})
}

func addValid(t *testing.T) {
	t.Parallel()
	globalT := t
	defer lifecycle.PurgeModels(globalT)

	_, s := testdata.NewAdminAuth(t)
	auth := httptests.NewRequestAuth(s)

	t.Run("parallel", func(t *testing.T) {
		testCases := []struct {
			description string
			code        int
			params      *articles.AddParams
		}{
			{
				"As few params as possible",
				http.StatusCreated,
				&articles.AddParams{Title: "My Super Article"},
			},
			{
				"All fields set",
				http.StatusCreated,
				&articles.AddParams{Title: "Title", Subtitle: "Subtitle", Description: "Description", Content: "Content"},
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()

				rec := callAdd(t, tc.params, auth)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusCreated {
					var a articles.Payload
					if err := json.NewDecoder(rec.Body).Decode(&a); err != nil {
						t.Fatal(err)
					}

					// Validate the article
					assert.NotEmpty(t, a.ID)
					assert.NotEmpty(t, a.Slug)
					assert.Nil(t, a.PublishedAt)

					// Validate the article's content
					assert.Equal(t, tc.params.Title, a.Title)
					assert.Equal(t, tc.params.Subtitle, a.Subtitle)
					assert.Equal(t, tc.params.Description, a.Description)
					assert.Equal(t, tc.params.Content, a.Content)

					lifecycle.SaveModels(globalT, &articles.Article{ID: a.ID})
				}
			})
		}
	})
}

// addAccessData tests that some wrong high level data are failing
func addAccessData(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	_, userSession := testdata.NewAuth(t)
	userAuth := httptests.NewRequestAuth(userSession)

	t.Run("parallel", func(t *testing.T) {
		testCases := []struct {
			description string
			code        int
			auth        *httptests.RequestAuth
			params      *articles.AddParams
		}{
			{
				"As anonymous",
				http.StatusUnauthorized,
				nil,
				&articles.AddParams{Title: "Title"},
			},
			{
				"As logged user",
				http.StatusUnauthorized,
				userAuth,
				&articles.AddParams{Title: "Title"},
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				rec := callAdd(t, tc.params, tc.auth)
				require.Equal(t, tc.code, rec.Code)
			})
		}
	})
}

func callAdd(t *testing.T, params *articles.AddParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: articles.Endpoints[articles.EndpointAdd],
		Params:   params,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}
