// +build integration

package organizations_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/primitives/models/lifecycle"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationDeleteHappyPath(t *testing.T) {
	dbCon := dependencies.DB

	defer lifecycle.PurgeModels(t, dbCon)
	_, admSession := testauth.NewAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)
	toDelete := testorganizations.NewOrganization(t, dbCon, nil)

	tests := []struct {
		description string
		code        int
		params      *organizations.DeleteParams
	}{
		{
			"Valid request should work",
			http.StatusNoContent,
			&organizations.DeleteParams{ID: toDelete.ID},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callDelete(t, tc.params, adminAuth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusNoContent {
				exists, err := organizations.Exists(dbCon, tc.params.ID)
				assert.NoError(t, err, "Exists() should have not failed")
				assert.False(t, exists, "the organization should no longer exists")
			}
		})
	}
}

func callDelete(t *testing.T, params *organizations.DeleteParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: organizations.Endpoints[organizations.EndpointDelete],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}
