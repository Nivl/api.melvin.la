package users_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"

	"github.com/Nivl/api.melvin.la/api/app/testhelpers"
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/components/users"
	"github.com/stretchr/testify/assert"
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
		params      *users.HandlerDeleteParams
		auth        *testhelpers.RequestAuth
	}{
		{
			"Not logged",
			http.StatusUnauthorized,
			&users.HandlerDeleteParams{ID: u1.ID.Hex()},
			nil,
		},
		{
			"Deleting an other user",
			http.StatusForbidden,
			&users.HandlerDeleteParams{ID: u1.ID.Hex()},
			testhelpers.NewRequestAuth(s2.ID, u2.ID),
		},
		{
			"Deleting without providing password",
			http.StatusUnauthorized,
			&users.HandlerDeleteParams{ID: u1.ID.Hex()},
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
		},
		// Keep this one last for u1 as it deletes the user
		{
			"Deleting user",
			http.StatusNoContent,
			&users.HandlerDeleteParams{ID: u1.ID.Hex(), CurrentPassword: "fake"},
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerDelete(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusNoContent {
				// We check that the user is still in DB but is flagged for deletion
				userID := bson.ObjectIdHex(tc.params.ID)
				var u auth.User
				if err := auth.QueryUsers().FindId(userID).One(&u); err != nil {
					t.Fatal(err)
				}

				if assert.NotNil(t, u, "User fully deleted") {
					assert.True(t, u.IsDeleted, "User not marked for deletion")
				}
			}
		})
	}
}

func callHandlerDelete(t *testing.T, params *users.HandlerDeleteParams, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: users.Endpoints[users.EndpointDelete],
		URI:      fmt.Sprintf("/users/%s", params.ID),
		Params:   params,
		Auth:     auth,
	}

	return testhelpers.NewRequest(ri)
}
