// +build integration

package education_test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-types/datetime"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/dchest/uniuri"
	"github.com/melvin-laplanche/ml-api/src/components/about/education"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationAdd(t *testing.T) {
	dbCon := deps.DB()

	defer lifecycle.PurgeModels(t, dbCon)
	_, admSession := testauth.NewAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)
	org := testorganizations.NewPersisted(t, dbCon, nil)
	deletedOrg := testorganizations.NewPersisted(t, dbCon, &organizations.Organization{
		DeletedAt: datetime.Now(),
	})

	tests := []struct {
		description string
		code        int
		params      *education.AddParams
	}{
		{
			"Valid Request should work",
			http.StatusCreated,
			&education.AddParams{
				OrganizationID: org.ID,
				Degree:         uniuri.New(),
				Location:       uniuri.New(),
				Description:    uniuri.New(),
				StartYear:      rand.Intn(100) + 1900,
			},
		},
		{
			"Using a trashed organisation should fail",
			http.StatusNotFound,
			&education.AddParams{
				OrganizationID: deletedOrg.ID,
				Degree:         uniuri.New(),
				Location:       uniuri.New(),
				Description:    uniuri.New(),
				StartYear:      rand.Intn(100) + 1900,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callAdd(t, tc.params, adminAuth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusCreated {
				ext := &education.Payload{}
				if err := json.NewDecoder(rec.Body).Decode(ext); err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, ext.ID)
				assert.NotNil(t, ext.CreatedAt)
				assert.NotNil(t, ext.UpdatedAt)
				assert.Nil(t, ext.DeletedAt)
				assert.Equal(t, tc.params.Degree, ext.Degree)
				assert.Equal(t, tc.params.Location, ext.Location)
				assert.Equal(t, tc.params.Description, ext.Description)
				assert.Equal(t, tc.params.StartYear, ext.StartYear)

				// clean the test
				extModel, err := education.GetByID(dbCon, ext.ID)
				if err != nil {
					t.Fatal(err)
				}
				extModel.Delete(dbCon)
			}
		})
	}
}

func callAdd(t *testing.T, params *education.AddParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: education.Endpoints[education.EndpointAdd],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}
