package users_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"

	"github.com/Nivl/api.melvin.la/api/app/testhelpers"
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/components/users"
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
		params      *users.HandlerDeleteParams
		auth        *testhelpers.RequestAuth
	}{
		{
			"Not logged",
			http.StatusUnauthorized,
			&users.HandlerDeleteParams{UUID: u1.UUID},
			nil,
		},
		{
			"Deleting an other user",
			http.StatusForbidden,
			&users.HandlerDeleteParams{UUID: u1.UUID},
			testhelpers.NewRequestAuth(s2.UUID, u2.UUID),
		},
		{
			"Deleting without providing password",
			http.StatusUnauthorized,
			&users.HandlerDeleteParams{UUID: u1.UUID},
			testhelpers.NewRequestAuth(s1.UUID, u1.UUID),
		},
		// Keep this one last for u1 as it deletes the user
		{
			"Deleting user",
			http.StatusNoContent,
			&users.HandlerDeleteParams{UUID: u1.UUID, CurrentPassword: "fake"},
			testhelpers.NewRequestAuth(s1.UUID, u1.UUID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerDelete(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusNoContent {
				// We check that the user is still in DB but is flagged for deletion
				var user auth.User
				stmt := "SELECT * FROM users WHERE uuid=$1 LIMIT 1"
				err := db.Get(&user, stmt, tc.params.UUID)
				if err != nil {
					t.Fatal(err)
				}

				if assert.NotEmpty(t, user.UUID, "User fully deleted") {
					assert.NotNil(t, user.DeletedAt, "User not marked for deletion")
				}
			}
		})
	}
}

func callHandlerDelete(t *testing.T, params *users.HandlerDeleteParams, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: users.Endpoints[users.EndpointDelete],
		URI:      fmt.Sprintf("/users/%s", params.UUID),
		Params:   params,
		Auth:     auth,
	}

	return testhelpers.NewRequest(ri)
}
