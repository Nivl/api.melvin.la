package sessions_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"

	"github.com/Nivl/api.melvin.la/api/app/testhelpers"
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/components/sessions"
	"github.com/stretchr/testify/assert"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
			&sessions.HandlerDeleteParams{ID: s1.ID.Hex()},
			nil,
		},
		{
			"Deleting an other user sessions",
			http.StatusNotFound,
			&sessions.HandlerDeleteParams{ID: s1.ID.Hex(), CurrentPassword: "fake"},
			testhelpers.NewRequestAuth(s2.ID, u2.ID),
		},
		{
			"Deleting an invalid ID",
			http.StatusNotFound,
			&sessions.HandlerDeleteParams{ID: "invalid", CurrentPassword: "fake"},
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
		},
		{
			"Deleting without providing password",
			http.StatusUnauthorized,
			&sessions.HandlerDeleteParams{ID: s1.ID.Hex()},
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
		},
		// Keep this one last for u1 as it deletes the session
		{
			"Deleting session",
			http.StatusNoContent,
			&sessions.HandlerDeleteParams{ID: s1.ID.Hex(), CurrentPassword: "fake"},
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerDelete(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusNoContent {
				// We check that the user is still in DB but is flagged for deletion
				sessionID := bson.ObjectIdHex(tc.params.ID)
				var s auth.Session
				err := auth.QuerySessions().FindId(sessionID).One(&s)
				if err != nil && err != mgo.ErrNotFound {
					t.Fatal(err)
				}

				if assert.NotEmpty(t, s.ID, "session fully deleted") {
					assert.True(t, s.IsDeleted, "User not marked for deletion")
				}
			}
		})
	}
}

func callHandlerDelete(t *testing.T, params *sessions.HandlerDeleteParams, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: sessions.Endpoints[sessions.EndpointDelete],
		URI:      fmt.Sprintf("/sessions/%s", params.ID),
		Params:   params,
		Auth:     auth,
	}

	return testhelpers.NewRequest(ri)
}
