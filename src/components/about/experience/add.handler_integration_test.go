// +build integration

package experience_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/testing/integration"
	"github.com/Nivl/go-types/date"
	"github.com/Nivl/go-types/datetime"
	"github.com/dchest/uniuri"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/melvin-laplanche/ml-api/src/components/api"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationAdd(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	_, admSession := testauth.NewPersistedAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)
	org := testorganizations.NewPersisted(t, dbCon, nil)
	deletedOrg := testorganizations.NewPersisted(t, dbCon, &organizations.Organization{
		DeletedAt: datetime.Now(),
	})

	tests := []struct {
		description string
		code        int
		params      *experience.AddParams
	}{
		{
			"Valid Request should work",
			http.StatusCreated,
			&experience.AddParams{
				OrganizationID: org.ID,
				JobTitle:       uniuri.New(),
				Location:       uniuri.New(),
				Description:    uniuri.New(),
				StartDate:      date.Today(),
			},
		},
		{
			"Using a trashed organisation should fail",
			http.StatusNotFound,
			&experience.AddParams{
				OrganizationID: deletedOrg.ID,
				JobTitle:       uniuri.New(),
				Location:       uniuri.New(),
				Description:    uniuri.New(),
				StartDate:      date.Today(),
			},
		},
	}

	t.Run("parallel", func(t *testing.T) {
		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				defer helper.RecoverPanic()

				rec := callAdd(t, tc.params, adminAuth, helper.Deps)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusCreated {
					ext := &experience.Payload{}
					if err := json.NewDecoder(rec.Body).Decode(ext); err != nil {
						t.Fatal(err)
					}

					assert.NotEmpty(t, ext.ID)
					assert.NotNil(t, ext.CreatedAt)
					assert.NotNil(t, ext.UpdatedAt)
					assert.Nil(t, ext.DeletedAt)
					assert.Equal(t, tc.params.JobTitle, ext.JobTitle)
					assert.Equal(t, tc.params.Location, *ext.Location)
					assert.Equal(t, tc.params.Description, *ext.Description)
					assert.Equal(t, tc.params.StartDate.String(), ext.StartDate.String())

					// clean the test
					extModel, err := experience.GetByID(dbCon, ext.ID)
					if err != nil {
						t.Fatal(err)
					}
					extModel.Delete(dbCon)
				}
			})
		}
	})
}

func callAdd(t *testing.T, params *experience.AddParams, auth *httptests.RequestAuth, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: experience.Endpoints[experience.EndpointAdd],
		Params:   params,
		Auth:     auth,
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
