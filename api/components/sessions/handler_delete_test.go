package sessions_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"

	"github.com/Nivl/api.melvin.la/api/app/testhelpers"
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/components/sessions"
	"github.com/Nivl/api.melvin.la/api/db"
	"github.com/stretchr/testify/assert"
)

func TestHandlerDelete(t *testing.T) {
	defer testhelpers.PurgeModels(t)

	u1, s1 := auth.NewTestAuth(t)
	u2, s2 := auth.NewTestAuth(t)
	testhelpers.SaveModel(t, u1)
	testhelpers.SaveModel(t, s1)
	testhelpers.SaveModel(t, u2)
	testhelpers.SaveModel(t, s2)

	tests := []struct {
		description string
		code        int
		params      *sessions.HandlerDeleteParams
		auth        *testhelpers.RequestAuth
	}{
		{
			"Not logged",
			http.StatusUnauthorized,
			&sessions.HandlerDeleteParams{Token: s1.UUID},
			nil,
		},
		{
			"Deleting an other user sessions",
			http.StatusNotFound,
			&sessions.HandlerDeleteParams{Token: s1.UUID, CurrentPassword: "fake"},
			testhelpers.NewRequestAuth(s2.UUID, u2.UUID),
		},
		{
			"Deleting an invalid ID",
			http.StatusNotFound,
			&sessions.HandlerDeleteParams{Token: "invalid", CurrentPassword: "fake"},
			testhelpers.NewRequestAuth(s1.UUID, u1.UUID),
		},
		{
			"Deleting without providing password",
			http.StatusUnauthorized,
			&sessions.HandlerDeleteParams{Token: s1.UUID},
			testhelpers.NewRequestAuth(s1.UUID, u1.UUID),
		},
		// Keep this one last for u1 as it deletes the session
		{
			"Deleting session",
			http.StatusNoContent,
			&sessions.HandlerDeleteParams{Token: s1.UUID, CurrentPassword: "fake"},
			testhelpers.NewRequestAuth(s1.UUID, u1.UUID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerDelete(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusNoContent {
				// We check that the user is still in DB but is flagged for deletion
				var session auth.Session
				stmt := "SELECT * FROM sessions WHERE uuid=$1 LIMIT 1"
				err := db.Get(&session, stmt, tc.params.Token)
				if err != nil {
					t.Fatal(err)
				}

				if assert.NotEmpty(t, session.UUID, "session fully deleted") {
					assert.NotNil(t, session.DeletedAt, "User not marked for deletion")
				}
			}
		})
	}
}

func callHandlerDelete(t *testing.T, params *sessions.HandlerDeleteParams, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: sessions.Endpoints[sessions.EndpointDelete],
		URI:      fmt.Sprintf("/sessions/%s", params.Token),
		Params:   params,
		Auth:     auth,
	}

	return testhelpers.NewRequest(ri)
}
