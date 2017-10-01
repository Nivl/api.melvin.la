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
	"github.com/Nivl/go-types/datetime"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience/testexperience"
	"github.com/melvin-laplanche/ml-api/src/components/api"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationGet(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	_, userSession := testauth.NewPersistedAuth(t, dbCon)
	_, adminSession := testauth.NewPersistedAdminAuth(t, dbCon)

	basicExp := testexperience.NewPersisted(t, dbCon, nil)

	orphanExp := testexperience.NewPersisted(t, dbCon, nil)
	orphanExp.Organization.DeletedAt = datetime.Now()
	orphanExp.Organization.Update(dbCon)

	tests := []struct {
		description  string
		code         int
		params       *experience.GetParams
		auth         *httptests.RequestAuth
		expectedData *experience.Experience
	}{
		{
			"Anonymous getting an orphan exp",
			http.StatusNotFound,
			&experience.GetParams{ID: orphanExp.ID},
			nil,
			orphanExp,
		},
		{
			"User getting an orphan exp",
			http.StatusNotFound,
			&experience.GetParams{ID: orphanExp.ID},
			httptests.NewRequestAuth(userSession),
			orphanExp,
		},
		{
			"Admin getting an orphan exp",
			http.StatusOK,
			&experience.GetParams{ID: orphanExp.ID},
			httptests.NewRequestAuth(adminSession),
			orphanExp,
		},
		{
			"Anonymous getting a basic exp",
			http.StatusOK,
			&experience.GetParams{ID: basicExp.ID},
			nil,
			basicExp,
		},
		{
			"Admin getting a basic exp",
			http.StatusOK,
			&experience.GetParams{ID: basicExp.ID},
			httptests.NewRequestAuth(adminSession),
			basicExp,
		},
	}

	t.Run("parallel", func(t *testing.T) {
		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				defer helper.RecoverPanic()

				rec := callGet(t, tc.params, tc.auth, helper.Deps)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusOK {
					var pld experience.Payload
					if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
						t.Fatal(err)
					}

					assert.Equal(t, tc.expectedData.ID, pld.ID, "ID should not have changed")
					assert.Equal(t, tc.expectedData.JobTitle, pld.JobTitle, "JobTitle should not have changed")
					assert.Equal(t, tc.expectedData.Location, pld.Location, "Location should not have changed")
					assert.Equal(t, tc.expectedData.Description, pld.Description, "Description should not have changed")
					assert.Equal(t, tc.expectedData.StartDate.String(), pld.StartDate.String(), "StartDate should not have changed")
					if tc.expectedData.EndDate != nil {
						assert.Equal(t, tc.expectedData.EndDate.String(), pld.EndDate.String(), "EndDate should not have changed")
					} else {
						assert.Nil(t, tc.expectedData.EndDate, "EndDate should still be nil")
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
	})
}

func callGet(t *testing.T, params *experience.GetParams, auth *httptests.RequestAuth, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: experience.Endpoints[experience.EndpointGet],
		Params:   params,
		Auth:     auth,
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
