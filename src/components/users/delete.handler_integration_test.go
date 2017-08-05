// +build integration

package users_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	defer lifecycle.PurgeModels(t, deps.DB)

	u1, s1 := testauth.NewAuth(t, deps.DB)
	_, s2 := testauth.NewAuth(t, deps.DB)

	tests := []struct {
		description string
		code        int
		params      *users.DeleteParams
		auth        *httptests.RequestAuth
	}{
		{
			"Not logged",
			http.StatusUnauthorized,
			&users.DeleteParams{ID: u1.ID, CurrentPassword: "psw"},
			nil,
		},
		{
			"Deleting an other user",
			http.StatusForbidden,
			&users.DeleteParams{ID: u1.ID, CurrentPassword: "psw"},
			httptests.NewRequestAuth(s2),
		},
		{
			"Deleting without providing password",
			http.StatusUnauthorized,
			&users.DeleteParams{ID: u1.ID, CurrentPassword: "psw"},
			httptests.NewRequestAuth(s1),
		},
		// Keep this one last for u1 as it deletes the user
		{
			"Deleting user",
			http.StatusNoContent,
			&users.DeleteParams{ID: u1.ID, CurrentPassword: "fake"},
			httptests.NewRequestAuth(s1),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callDelete(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusNoContent {
				// We check that the user has been deleted
				var user auth.User
				stmt := "SELECT * FROM users WHERE id=$1 LIMIT 1"
				err := db.Get(deps.DB, &user, stmt, tc.params.ID)
				if err != nil {
					t.Fatal(err)
				}

				assert.Empty(t, user.ID, "User not deleted")
			}
		})
	}
}

func callDelete(t *testing.T, params *users.DeleteParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: users.Endpoints[users.EndpointDelete],
		Params:   params,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}
