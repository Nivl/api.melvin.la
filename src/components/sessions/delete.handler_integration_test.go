// +build integration

package sessions_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/melvin-laplanche/ml-api/src/components/sessions"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	dbCon := deps.DB()
	defer lifecycle.PurgeModels(t, dbCon)

	// Do not delete safeSession
	u1, safeSession := testauth.NewAuth(t, dbCon)

	// We create a couple of sessions for the same user
	toDelete2 := testauth.NewPersistedSession(t, dbCon, u1)
	toDelete3 := testauth.NewPersistedSession(t, dbCon, u1)

	// We create a other session attached to an other user
	_, randomSession := testauth.NewAuth(t, dbCon)

	tests := []struct {
		description string
		code        int
		params      *sessions.DeleteParams
		auth        *httptests.RequestAuth
	}{
		{
			"Not logged",
			http.StatusUnauthorized,
			&sessions.DeleteParams{Token: safeSession.ID},
			nil,
		},
		{
			"Deleting an other user sessions",
			http.StatusNotFound,
			&sessions.DeleteParams{Token: safeSession.ID, CurrentPassword: "fake"},
			httptests.NewRequestAuth(randomSession),
		},
		{
			"Deleting an invalid ID",
			http.StatusBadRequest,
			&sessions.DeleteParams{Token: "invalid", CurrentPassword: "fake"},
			httptests.NewRequestAuth(safeSession),
		},
		{
			"Deleting a different session without providing password",
			http.StatusUnauthorized,
			&sessions.DeleteParams{Token: toDelete2.ID},
			httptests.NewRequestAuth(safeSession),
		},
		{
			"Deleting a different session",
			http.StatusNoContent,
			&sessions.DeleteParams{Token: toDelete2.ID, CurrentPassword: "fake"},
			httptests.NewRequestAuth(safeSession),
		},
		{
			"Deleting current session",
			http.StatusNoContent,
			&sessions.DeleteParams{Token: toDelete3.ID},
			httptests.NewRequestAuth(toDelete3),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callDelete(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusNoContent {
				// We check that the user is still in DB but is flagged for deletion
				var session auth.Session
				stmt := "SELECT * FROM sessions WHERE id=$1 LIMIT 1"
				err := dbCon.Get(&session, stmt, tc.params.Token)
				assert.Equal(t, sql.ErrNoRows, err, "session should be deleted")
			}
		})
	}
}

func callDelete(t *testing.T, params *sessions.DeleteParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: sessions.Endpoints[sessions.EndpointDelete],
		Params:   params,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}
