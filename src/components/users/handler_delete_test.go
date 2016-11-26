package users_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"

	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/auth/authtest"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/melvin-laplanche/ml-api/src/db"
	"github.com/melvin-laplanche/ml-api/src/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestHandlerDelete(t *testing.T) {
	defer testhelpers.PurgeModels(t)

	u1, s1 := authtest.NewAuth(t)
	u2, s2 := authtest.NewAuth(t)

	tests := []struct {
		description string
		code        int
		params      *users.HandlerDeleteParams
		auth        *testhelpers.RequestAuth
	}{
		{
			"Not logged",
			http.StatusUnauthorized,
			&users.HandlerDeleteParams{ID: u1.ID},
			nil,
		},
		{
			"Deleting an other user",
			http.StatusForbidden,
			&users.HandlerDeleteParams{ID: u1.ID},
			testhelpers.NewRequestAuth(s2.ID, u2.ID),
		},
		{
			"Deleting without providing password",
			http.StatusUnauthorized,
			&users.HandlerDeleteParams{ID: u1.ID},
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
		},
		// Keep this one last for u1 as it deletes the user
		{
			"Deleting user",
			http.StatusNoContent,
			&users.HandlerDeleteParams{ID: u1.ID, CurrentPassword: "fake"},
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerDelete(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusNoContent {
				// We check that the user is still in DB but is flagged for deletion
				var user auth.User
				stmt := "SELECT * FROM users WHERE id=$1 LIMIT 1"
				err := db.Get(&user, stmt, tc.params.ID)
				if err != nil {
					t.Fatal(err)
				}

				if assert.NotEmpty(t, user.ID, "User fully deleted") {
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
		URI:      fmt.Sprintf("/users/%s", params.ID),
		Params:   params,
		Auth:     auth,
	}

	return testhelpers.NewRequest(ri)
}
