package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/api.melvin.la/api/app/testhelpers"
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/router"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestEndpointExecution test that an endpoint gets executed with no issue.
// It only tests that the main middleware does what it is supposed to do, and
// therefore does not tests details (like checking the params are parsed correctly)
func TestEndpointExecution(t *testing.T) {
	// Handler used for our request. We just need to know if it is called or not
	hdlr := func(req *router.Request) {
		req.NoContent()
	}

	u, s := auth.NewTestAuth(t)
	testhelpers.SaveModels(t, u, s)
	defer testhelpers.PurgeModels(t)

	tests := []struct {
		description string
		endpoint    *router.Endpoint
		url         string
		auth        *testhelpers.RequestAuth
		code        int
	}{
		{
			"Basic public GET",
			&router.Endpoint{Verb: "GET", Path: "/items", Handler: hdlr},
			"/items",
			nil,
			http.StatusNoContent,
		},
		{
			"Private GET as anonymous",
			&router.Endpoint{Verb: "GET", Path: "/items/{id}", Handler: hdlr, Auth: router.LoggedUser},
			"/items/item-id",
			nil,
			http.StatusUnauthorized,
		},
		{
			"Private GET as logged user",
			&router.Endpoint{Verb: "GET", Path: "/items/{id}", Handler: hdlr, Auth: router.LoggedUser},
			"/items/item-id",
			testhelpers.NewRequestAuth(s.ID, u.ID),
			http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := execHandler(t, tc.endpoint, tc.url, tc.auth)
			assert.Equal(t, tc.code, rec.Code)
		})
	}
}

func execHandler(t *testing.T, e *router.Endpoint, url string, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
	r := mux.NewRouter()
	r.Methods(e.Verb).Path(e.Path).Handler(router.Handler(e))

	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: e,
		URI:      url,
		Router:   r,
		Auth:     auth,
	}

	return testhelpers.NewRequest(ri)
}
