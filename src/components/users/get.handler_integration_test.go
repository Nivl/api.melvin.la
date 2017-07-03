// +build integration

package users_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/primitives/models/lifecycle"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	defer lifecycle.PurgeModels(t, deps.DB)

	u1, s1 := testauth.NewAuth(t, deps.DB)
	_, s2 := testauth.NewAuth(t, deps.DB)

	tests := []struct {
		description string
		code        int
		params      *users.GetParams
		auth        *httptests.RequestAuth
	}{
		{
			"Not logged",
			http.StatusOK,
			&users.GetParams{ID: u1.ID},
			nil,
		},
		{
			"Getting an other user",
			http.StatusOK,
			&users.GetParams{ID: u1.ID},
			httptests.NewRequestAuth(s2),
		},
		{
			"Getting own data",
			http.StatusOK,
			&users.GetParams{ID: u1.ID},
			httptests.NewRequestAuth(s1),
		},
		{
			"Getting un-existing user with valid ID",
			http.StatusNotFound,
			&users.GetParams{ID: "f76700e7-988c-4ae9-9f02-ac3f9d7cd88e"},
			nil,
		},
		{
			"Getting un-existing user with invalid ID",
			http.StatusBadRequest,
			&users.GetParams{ID: "invalidID"},
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callGet(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusOK {
				var u users.Payload
				if err := json.NewDecoder(rec.Body).Decode(&u); err != nil {
					t.Fatal(err)
				}

				if assert.Equal(t, tc.params.ID, u.ID, "Not the same user") {
					// User access their own data
					if tc.auth != nil && u.ID == tc.auth.UserID {
						assert.NotEmpty(t, u.Email, "Same user needs their private data")
					} else { // user access an other user data
						assert.Empty(t, u.Email, "Should not return private data")
					}
				}
			}
		})
	}
}

func callGet(t *testing.T, params *users.GetParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: users.Endpoints[users.EndpointGet],
		Params:   params,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}
