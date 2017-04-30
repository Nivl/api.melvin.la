package users_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/primitives/models/lifecycle"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	globalT := t
	defer lifecycle.PurgeModels(t)

	tests := []struct {
		description string
		code        int
		params      *users.AddParams
	}{
		{
			"Empty User",
			http.StatusBadRequest,
			&users.AddParams{},
		},
		{
			"Valid User",
			http.StatusCreated,
			&users.AddParams{Name: "Name", Email: "email+TestAdd@fake.com", Password: "password"},
		},
		{
			"Duplicate Email",
			http.StatusConflict,
			&users.AddParams{Name: "Name", Email: "email+TestAdd@fake.com", Password: "password"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerAdd(t, tc.params)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusCreated {
				var u users.Payload
				if err := json.NewDecoder(rec.Body).Decode(&u); err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, u.ID)
				assert.Equal(t, tc.params.Email, u.Email)
				lifecycle.SaveModels(globalT, &auth.User{ID: u.ID})
			}
		})
	}
}

func callHandlerAdd(t *testing.T, params *users.AddParams) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: users.Endpoints[users.EndpointAdd],
		Params:   params,
	}

	return httptests.NewRequest(t, ri)
}
