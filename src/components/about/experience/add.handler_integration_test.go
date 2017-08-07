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
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/dchest/uniuri"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationAdd(t *testing.T) {
	dbCon := dependencies.DB

	defer lifecycle.PurgeModels(t, dbCon)
	_, admSession := testauth.NewAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)
	org := testorganizations.NewOrganization(t, dbCon, nil)
	deletedOrg := testorganizations.NewOrganization(t, dbCon, &organizations.Organization{
		DeletedAt: db.Now(),
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
				StartDate:      db.Today(),
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
				StartDate:      db.Today(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callAdd(t, tc.params, adminAuth)
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
				assert.Equal(t, tc.params.Location, ext.Location)
				assert.Equal(t, tc.params.Description, ext.Description)
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
}

func callAdd(t *testing.T, params *experience.AddParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: experience.Endpoints[experience.EndpointAdd],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}
