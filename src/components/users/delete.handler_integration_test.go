// +build integration

package users_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/testing/integration"
	"github.com/melvin-laplanche/ml-api/src/components/api"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	u1, s1 := testauth.NewPersistedAuth(t, dbCon)
	_, s2 := testauth.NewPersistedAuth(t, dbCon)

	userToDelete, sessionToDelete := testauth.NewPersistedAuth(t, dbCon)

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
		{
			"Deleting user",
			http.StatusNoContent,
			&users.DeleteParams{ID: userToDelete.ID, CurrentPassword: "fake"},
			httptests.NewRequestAuth(sessionToDelete),
		},
	}

	t.Run("parallel", func(t *testing.T) {
		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				defer helper.RecoverPanic()

				rec := callDelete(t, tc.params, tc.auth, helper.Deps)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusNoContent {
					// We check that the user has been deleted
					var user auth.User
					stmt := "SELECT * FROM users WHERE id=$1 LIMIT 1"
					err := dbCon.Get(&user, stmt, tc.params.ID)
					assert.Equal(t, sql.ErrNoRows, err, "User not deleted")
				}
			})
		}
	})
}

func callDelete(t *testing.T, params *users.DeleteParams, auth *httptests.RequestAuth, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: users.Endpoints[users.EndpointDelete],
		Params:   params,
		Auth:     auth,
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
