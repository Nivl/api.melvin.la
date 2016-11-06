package articles_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/api.melvin.la/api/app/testhelpers"
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/components/blog/articles"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestHandlerAdd(t *testing.T) {
	globalT := t
	defer testhelpers.PurgeModels(t)

	u1, s1 := auth.NewTestAuth(t)
	testhelpers.SaveModel(t, u1)
	testhelpers.SaveModel(t, s1)

	tests := []struct {
		description string
		code        int
		params      *articles.HandlerAddParams
		auth        *testhelpers.RequestAuth
	}{
		{
			"No logged",
			http.StatusBadRequest,
			&articles.HandlerAddParams{},
			nil,
		},
		{
			"No Title",
			http.StatusBadRequest,
			&articles.HandlerAddParams{},
			testhelpers.NewRequestAuth(s1.UUID, u1.UUID),
		},
		{
			"Title filled with spaces",
			http.StatusBadRequest,
			&articles.HandlerAddParams{Title: "       "},
			testhelpers.NewRequestAuth(s1.UUID, u1.UUID),
		},
		{
			"As few params as possible",
			http.StatusCreated,
			&articles.HandlerAddParams{Title: "My Super Article"},
			testhelpers.NewRequestAuth(s1.UUID, u1.UUID),
		},
		{
			"Duplicate title",
			http.StatusCreated,
			&articles.HandlerAddParams{Title: "My Super Article"},
			testhelpers.NewRequestAuth(s1.UUID, u1.UUID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerAdd(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusCreated {
				var a articles.Exportable
				if err := json.NewDecoder(rec.Body).Decode(&a); err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, a.ID)
				assert.NotEmpty(t, a.Slug)
				assert.Equal(t, tc.params.Title, a.Title)
				testhelpers.SaveModel(globalT, &articles.Article{ID: bson.ObjectIdHex(a.ID)})
			}
		})
	}
}

func callHandlerAdd(t *testing.T, params *articles.HandlerAddParams, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: articles.Endpoints[articles.EndpointAdd],
		URI:      "/blog/articles/",
		Params:   params,
		Auth:     auth,
	}

	return testhelpers.NewRequest(ri)
}
