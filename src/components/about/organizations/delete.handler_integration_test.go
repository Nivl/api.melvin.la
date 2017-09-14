// +build integration

package organizations_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-types/datetime"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationDeleteHappyPath(t *testing.T) {
	dbCon := deps.DB()

	defer lifecycle.PurgeModels(t, dbCon)
	_, admSession := testauth.NewAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)
	basicOrg := testorganizations.NewPersisted(t, dbCon, nil)
	trashedOrg := testorganizations.NewPersisted(t, dbCon, &organizations.Organization{
		DeletedAt: datetime.Now(),
	})

	tests := []struct {
		description string
		code        int
		params      *organizations.DeleteParams
	}{
		{
			"Valid request should work",
			http.StatusNoContent,
			&organizations.DeleteParams{ID: basicOrg.ID},
		},
		{
			"trashed org should work",
			http.StatusNoContent,
			&organizations.DeleteParams{ID: trashedOrg.ID},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callDelete(t, tc.params, adminAuth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusNoContent {
				_, err := organizations.GetAnyByID(dbCon, tc.params.ID)
				assert.True(t, apierror.IsNotFound(err), "GetByID() should have failed with an IsNotFound error")
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
