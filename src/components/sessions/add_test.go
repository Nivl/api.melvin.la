package sessions_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/primitives/models/lifecycle"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/security/auth/testdata"
	"github.com/melvin-laplanche/ml-api/src/components/sessions"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	defer lifecycle.PurgeModels(t)
	u1 := testdata.NewUser(t, nil)

	tests := []struct {
		description string
		code        int
		params      *sessions.AddParams
	}{
		{
			"Invalid email",
			http.StatusBadRequest,
			&sessions.AddParams{Email: "invalid@fake.com", Password: "fake"},
		},
		{
			"Invalid password",
			http.StatusBadRequest,
			&sessions.AddParams{Email: u1.Email, Password: "invalid"},
		},
		{
			"Valid Request",
			http.StatusCreated,
			&sessions.AddParams{Email: u1.Email, Password: "fake"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callAdd(t, tc.params)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusCreated {
				var session sessions.Payload
				if err := json.NewDecoder(rec.Body).Decode(&session); err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, session.Token)
				assert.Equal(t, u1.ID, session.UserID)

				// clean the test
				(&auth.Session{ID: session.Token}).FullyDelete()
			}
		})
	}
}

func callAdd(t *testing.T, params *sessions.AddParams) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: sessions.Endpoints[sessions.EndpointAdd],
		Params:   params,
	}

	return httptests.NewRequest(t, ri)
}
