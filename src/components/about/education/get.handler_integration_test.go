// +build integration

package education_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/types/datetime"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/melvin-laplanche/ml-api/src/components/about/education"
	"github.com/melvin-laplanche/ml-api/src/components/about/education/testeducation"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationGet(t *testing.T) {
	dbCon := dependencies.DB
	defer lifecycle.PurgeModels(t, dbCon)

	_, userSession := testauth.NewAuth(t, dbCon)
	_, adminSession := testauth.NewAdminAuth(t, dbCon)

	basicExp := testeducation.NewPersisted(t, dbCon, nil)

	orphanEdu := testeducation.NewPersisted(t, dbCon, nil)
	orphanEdu.Organization.DeletedAt = datetime.Now()
	orphanEdu.Organization.Update(dbCon)

	tests := []struct {
		description  string
		code         int
		params       *education.GetParams
		auth         *httptests.RequestAuth
		expectedData *education.Education
	}{
		{
			"Anonymous getting an orphan edu",
			http.StatusNotFound,
			&education.GetParams{ID: orphanEdu.ID},
			nil,
			orphanEdu,
		},
		{
			"User getting an orphan edu",
			http.StatusNotFound,
			&education.GetParams{ID: orphanEdu.ID},
			httptests.NewRequestAuth(userSession),
			orphanEdu,
		},
		{
			"Admin getting an orphan edu",
			http.StatusOK,
			&education.GetParams{ID: orphanEdu.ID},
			httptests.NewRequestAuth(adminSession),
			orphanEdu,
		},
		{
			"Anonymous getting a basic edu",
			http.StatusOK,
			&education.GetParams{ID: basicExp.ID},
			nil,
			basicExp,
		},
		{
			"Admin getting a basic edu",
			http.StatusOK,
			&education.GetParams{ID: basicExp.ID},
			httptests.NewRequestAuth(adminSession),
			basicExp,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callGet(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusOK {
				var pld education.Payload
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tc.expectedData.ID, pld.ID, "ID should not have changed")
				assert.Equal(t, tc.expectedData.Degree, pld.Degree, "Degree should not have changed")
				assert.Equal(t, tc.expectedData.Location, pld.Location, "Location should not have changed")
				assert.Equal(t, tc.expectedData.Description, pld.Description, "Description should not have changed")
				assert.Equal(t, tc.expectedData.StartYear, pld.StartYear, "StartYear should not have changed")
				if tc.expectedData.EndYear != nil {
					assert.Equal(t, tc.expectedData.EndYear, pld.EndYear, "EndYear should not have changed")
				} else {
					assert.Nil(t, tc.expectedData.EndYear, "EndYear should still be nil")
				}

				// we assume that if we have one field, we have all the fields
				assert.NotNil(t, pld.Organization, "Organization should not have been nil")
				assert.Equal(t, tc.expectedData.Organization.ID, pld.Organization.ID, "Organization.ID should not have changed")

				if tc.auth != nil && tc.auth.SessionID == adminSession.ID {
					assert.NotNil(t, pld.CreatedAt, "CreatedAt should have been set")
					assert.NotNil(t, pld.UpdatedAt, "UpdatedAt should have been set")
					assert.NotNil(t, pld.Organization.CreatedAt, "Organization.CreatedAt should have been set")
					assert.NotNil(t, pld.Organization.UpdatedAt, "Organization.UpdatedAt should have been set")

					if tc.expectedData.DeletedAt != nil {
						assert.NotNil(t, pld.DeletedAt, "DeletedAt should have been set")
					}

					if tc.expectedData.Organization.DeletedAt != nil {
						assert.NotNil(t, pld.Organization.DeletedAt, "Organization.DeletedAt should have been set")
					}
				} else {
					assert.Nil(t, pld.CreatedAt, "CreatedAt should have not been set")
					assert.Nil(t, pld.UpdatedAt, "UpdatedAt should have not been set")
					assert.Nil(t, pld.DeletedAt, "DeletedAt should have not been set")
					assert.Nil(t, pld.Organization.CreatedAt, "Organization.CreatedAt should have not been set")
					assert.Nil(t, pld.Organization.UpdatedAt, "Organization.UpdatedAt should have not been set")
					assert.Nil(t, pld.Organization.DeletedAt, "Organization.DeletedAt should have not been set")
				}
			}
		})
	}
}

func callGet(t *testing.T, params *education.GetParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: education.Endpoints[education.EndpointGet],
		Params:   params,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}
