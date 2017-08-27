// +build integration

package users_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/Nivl/go-rest-tools/types/ptrs"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/melvin-laplanche/ml-api/src/components/users/testusers"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	defer lifecycle.PurgeModels(t, deps.DB)

	u1, s1 := testusers.NewAuth(t, deps.DB)
	u2, s2 := testusers.NewAuth(t, deps.DB)

	tests := []struct {
		description string
		code        int
		params      *users.UpdateParams
		auth        *httptests.RequestAuth
	}{
		{
			"Not logged",
			http.StatusUnauthorized,
			&users.UpdateParams{ID: u1.ID},
			nil,
		},
		{
			"Updating an other user",
			http.StatusForbidden,
			&users.UpdateParams{ID: u1.ID},
			httptests.NewRequestAuth(s2),
		},
		{
			"Updating email without providing password",
			http.StatusUnauthorized,
			&users.UpdateParams{ID: u1.ID, Email: "melvin@fake.io"},
			httptests.NewRequestAuth(s1),
		},
		{
			"Updating password without providing current Password",
			http.StatusUnauthorized,
			&users.UpdateParams{ID: u1.ID, NewPassword: "TestUpdateUser"},
			httptests.NewRequestAuth(s1),
		},
		{
			"Updating regular field",
			http.StatusOK,
			&users.UpdateParams{ID: u1.ID, Name: "Melvin"},
			httptests.NewRequestAuth(s1),
		},
		{
			"Updating email to a used one",
			http.StatusConflict,
			&users.UpdateParams{ID: u1.ID, CurrentPassword: "fake", Email: u2.Email},
			httptests.NewRequestAuth(s1),
		},
		{
			"Updating facebook",
			http.StatusOK,
			&users.UpdateParams{ID: u1.ID, CurrentPassword: "fake", FacebookUsername: ptrs.NewString("fb")},
			httptests.NewRequestAuth(s1),
		},
		// Keep this one last for u1 as it changes the password
		{
			"Updating password",
			http.StatusOK,
			&users.UpdateParams{ID: u1.ID, CurrentPassword: "fake", NewPassword: "TestUpdateUser"},
			httptests.NewRequestAuth(s1),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callUpdate(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusOK {
				var u users.ProfilePayload
				if err := json.NewDecoder(rec.Body).Decode(&u); err != nil {
					t.Fatal(err)
				}

				if tc.params.Name != "" {
					assert.Equal(t, tc.params.Name, u.Name)
				}
				if tc.params.Email != "" {
					assert.Equal(t, tc.params.Email, u.Email)
				}
				if tc.params.FacebookUsername != nil {
					assert.Equal(t, *tc.params.FacebookUsername, u.FacebookUsername)
				}

				if tc.params.NewPassword != "" {
					// To check the password has been updated with need to get the
					// encrypted version, and compare it to the raw one
					updatedUser, err := auth.GetUserByID(deps.DB, tc.params.ID)
					if err != nil {
						t.Fatal(err)
					}

					hash := updatedUser.Password
					assert.True(t, auth.IsPasswordValid(hash, tc.params.NewPassword))
				}
			}
		})
	}
}

func callUpdate(t *testing.T, params *users.UpdateParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: users.Endpoints[users.EndpointUpdate],
		Params:   params,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}
