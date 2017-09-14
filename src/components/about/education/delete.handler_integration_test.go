// +build integration

package education_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-types/datetime"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/melvin-laplanche/ml-api/src/components/about/education"
	"github.com/melvin-laplanche/ml-api/src/components/about/education/testeducation"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationDeleteHappyPath(t *testing.T) {
	dbCon := deps.DB()

	defer lifecycle.PurgeModels(t, dbCon)
	_, admSession := testauth.NewPersistedAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)
	basicExp := testeducation.NewPersisted(t, dbCon, nil)
	trashedExp := testeducation.NewPersisted(t, dbCon, &education.Education{
		DeletedAt: datetime.Now(),
	})

	tests := []struct {
		description string
		code        int
		params      *education.DeleteParams
	}{
		{
			"Valid request should work",
			http.StatusNoContent,
			&education.DeleteParams{ID: basicExp.ID},
		},
		{
			"trashed exp should work",
			http.StatusNoContent,
			&education.DeleteParams{ID: trashedExp.ID},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callDelete(t, tc.params, adminAuth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusNoContent {
				_, err := education.GetAnyByID(dbCon, tc.params.ID)
				assert.True(t, apierror.IsNotFound(err), "GetByID() should have failed with an IsNotFound error")
			}
		})
	}
}

func callDelete(t *testing.T, params *education.DeleteParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: education.Endpoints[education.EndpointDelete],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}
