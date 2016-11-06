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

func TestHandlerUpdate(t *testing.T) {
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
		params      *users.HandlerUpdateParams
		auth        *testhelpers.RequestAuth
	}{
		{
			"Not logged",
			http.StatusUnauthorized,
			&users.HandlerUpdateParams{UUID: u1.UUID},
			nil,
		},
		{
			"Updating an other user",
			http.StatusForbidden,
			&users.HandlerUpdateParams{UUID: u1.UUID},
			testhelpers.NewRequestAuth(s2.UUID, u2.UUID),
		},
		{
			"Updating email without providing password",
			http.StatusUnauthorized,
			&users.HandlerUpdateParams{UUID: u1.UUID, Email: "melvin@fake.io"},
			testhelpers.NewRequestAuth(s1.UUID, u1.UUID),
		},
		{
			"Updating password without providing current Password",
			http.StatusUnauthorized,
			&users.HandlerUpdateParams{UUID: u1.UUID, NewPassword: "TestUpdateUser"},
			testhelpers.NewRequestAuth(s1.UUID, u1.UUID),
		},
		{
			"Updating regular field",
			http.StatusOK,
			&users.HandlerUpdateParams{UUID: u1.UUID, Name: "Melvin"},
			testhelpers.NewRequestAuth(s1.UUID, u1.UUID),
		},
		{
			"Updating email to a used one",
			http.StatusConflict,
			&users.HandlerUpdateParams{UUID: u1.UUID, CurrentPassword: "fake", Email: u2.Email},
			testhelpers.NewRequestAuth(s1.UUID, u1.UUID),
		},
		// Keep this one last for u1 as it changes the password
		{
			"Updating password",
			http.StatusOK,
			&users.HandlerUpdateParams{UUID: u1.UUID, CurrentPassword: "fake", NewPassword: "TestUpdateUser"},
			testhelpers.NewRequestAuth(s1.UUID, u1.UUID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerUpdate(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusOK {
				var u users.PrivatePayload
				if err := json.NewDecoder(rec.Body).Decode(&u); err != nil {
					t.Fatal(err)
				}

				if tc.params.Name != "" {
					assert.NotEmpty(t, tc.params.Name, u.Name)
				}

				if tc.params.Email != "" {
					assert.NotEmpty(t, tc.params.Email, u.Email)
				}

				if tc.params.NewPassword != "" {
					// To check the password has been updated with need to get the
					// encrypted version, and compare it to the raw one
					updatedUser, err := auth.GetUser(tc.params.UUID)
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

func callHandlerUpdate(t *testing.T, params *users.HandlerUpdateParams, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: users.Endpoints[users.EndpointUpdate],
		URI:      fmt.Sprintf("/users/%s", params.UUID),
		Params:   params,
		Auth:     auth,
	}

	return testhelpers.NewRequest(ri)
}
