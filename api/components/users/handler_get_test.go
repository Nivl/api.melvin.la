package users_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"

	"github.com/Nivl/api.melvin.la/api/app/testhelpers"
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/components/users"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGet(t *testing.T) {
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
		params      *users.HandlerGetParams
		auth        *testhelpers.RequestAuth
	}{
		{
			"Not logged",
			http.StatusOK,
			&users.HandlerGetParams{ID: u1.ID},
			nil,
		},
		{
			"Getting an other user",
			http.StatusOK,
			&users.HandlerGetParams{ID: u1.ID},
			testhelpers.NewRequestAuth(s2.ID, u2.ID),
		},
		{
			"Getting own data",
			http.StatusOK,
			&users.HandlerGetParams{ID: u1.ID},
			testhelpers.NewRequestAuth(s1.ID, u1.ID),
		},
		{
			"Getting un-existing user with valid ID",
			http.StatusNotFound,
			&users.HandlerGetParams{ID: "f76700e7-988c-4ae9-9f02-ac3f9d7cd88e"},
			nil,
		},
		{
			"Getting un-existing user with invalid ID",
			http.StatusBadRequest,
			&users.HandlerGetParams{ID: "invalidID"},
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerGet(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusOK {
				var u users.PrivatePayload
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

func callHandlerGet(t *testing.T, params *users.HandlerGetParams, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: users.Endpoints[users.EndpointGet],
		URI:      fmt.Sprintf("/users/%s", params.ID),
		Params:   params,
		Auth:     auth,
	}

	return testhelpers.NewRequest(ri)
}
