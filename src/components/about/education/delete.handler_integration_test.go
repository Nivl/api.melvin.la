// +build integration

package education_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/testing/integration"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-types/datetime"
	"github.com/melvin-laplanche/ml-api/src/components/about/education"
	"github.com/melvin-laplanche/ml-api/src/components/about/education/testeducation"
	"github.com/melvin-laplanche/ml-api/src/components/api"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationDeleteHappyPath(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

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

	t.Run("parallel", func(t *testing.T) {
		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				defer helper.RecoverPanic()

				rec := callDelete(t, tc.params, adminAuth, helper.Deps)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusNoContent {
					_, err := education.GetAnyByID(dbCon, tc.params.ID)
					assert.True(t, apierror.IsNotFound(err), "GetByID() should have failed with an IsNotFound error")
				}
			})
		}
	})
}

func callDelete(t *testing.T, params *education.DeleteParams, auth *httptests.RequestAuth, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: education.Endpoints[education.EndpointDelete],
		Params:   params,
		Auth:     auth,
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
