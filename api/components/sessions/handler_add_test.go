package sessions_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/api.melvin.la/api/app/testhelpers"
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/components/sessions"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestHandlerAdd(t *testing.T) {
	defer testhelpers.PurgeModels(t)

	u1 := auth.NewTestUser(t, nil)
	testhelpers.SaveModel(t, u1)

	tests := []struct {
		description string
		code        int
		params      *sessions.HandlerAddParams
	}{
		{
			"Invalid email",
			http.StatusBadRequest,
			&sessions.HandlerAddParams{Email: "invalid@fake.com", Password: "fake"},
		},
		{
			"Invalid password",
			http.StatusBadRequest,
			&sessions.HandlerAddParams{Email: u1.Email, Password: "invalid"},
		},
		{
			"Valid Request",
			http.StatusCreated,
			&sessions.HandlerAddParams{Email: u1.Email, Password: "fake"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerAdd(t, tc.params)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusCreated {
				var session sessions.Payload
				if err := json.NewDecoder(rec.Body).Decode(&session); err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, session.Token)
				assert.Equal(t, u1.ID.Hex(), session.UserID)

				// clean the test
				(&auth.Session{ID: bson.ObjectIdHex(session.Token)}).FullyDelete()
			}
		})
	}
}

func callHandlerAdd(t *testing.T, params *sessions.HandlerAddParams) *httptest.ResponseRecorder {
	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: sessions.Endpoints[sessions.EndpointAdd],
		URI:      "/sessions/",
		Params:   params,
	}

	return testhelpers.NewRequest(ri)
}
