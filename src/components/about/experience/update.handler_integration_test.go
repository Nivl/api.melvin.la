// +build integration

package experience_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/storage/db"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/Nivl/go-rest-tools/types/ptrs"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience/testexperience"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationUpdate(t *testing.T) {
	dbCon := dependencies.DB

	defer lifecycle.PurgeModels(t, dbCon)
	_, admSession := testauth.NewAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)

	noop := testexperience.NewPersisted(t, dbCon, nil)
	changeAll := testexperience.NewPersisted(t, dbCon, nil)
	toUntrash := testexperience.NewPersisted(t, dbCon, &experience.Experience{
		DeletedAt: db.Now(),
	})

	tests := []struct {
		description string
		code        int
		toUpdate    *experience.Experience
		params      *experience.UpdateParams
	}{
		{
			"Valid request should work",
			http.StatusOK,
			changeAll,
			&experience.UpdateParams{
				ID:          changeAll.ID,
				JobTitle:    ptrs.NewString("job title"),
				Description: ptrs.NewString("description"),
				Location:    ptrs.NewString("Location"),
				StartDate:   db.Today(),
				EndDate:     db.Today(),
				InTrash:     ptrs.NewBool(true),
			},
		},
		{
			"Untrash should work",
			http.StatusOK,
			toUntrash,
			&experience.UpdateParams{
				ID:      toUntrash.ID,
				InTrash: ptrs.NewBool(false),
			},
		},
		{
			"Noop should work",
			http.StatusOK,
			noop,
			&experience.UpdateParams{
				ID: noop.ID,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callUpdate(t, tc.params, adminAuth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusOK {
				var pld *experience.Payload
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tc.toUpdate.ID, pld.ID, "ID should have not changed")
				if tc.params.JobTitle != nil {
					assert.Equal(t, *tc.params.JobTitle, pld.JobTitle, "JobTitle should have changed")
				} else {
					assert.Equal(t, tc.toUpdate.JobTitle, pld.JobTitle, "JobTitle should have not changed")
				}

				if tc.params.Location != nil {
					assert.Equal(t, *tc.params.Location, pld.Location, "Location should have changed")
				} else {
					assert.Equal(t, tc.toUpdate.Location, pld.Location, "Location should have not changed")
				}

				if tc.params.Description != nil {
					assert.Equal(t, *tc.params.Description, pld.Description, "Description should have changed")
				} else {
					assert.Equal(t, tc.toUpdate.Description, pld.Description, "Description should have not changed")
				}

				if tc.params.StartDate != nil {
					assert.Equal(t, tc.params.StartDate.String(), pld.StartDate.String(), "StartDate should have changed")
				} else {
					assert.Equal(t, tc.toUpdate.StartDate.String(), pld.StartDate.String(), "StartDate should have not changed")
				}

				if tc.params.EndDate != nil {
					assert.Equal(t, tc.params.EndDate.String(), pld.EndDate.String(), "EndDate should have changed")
				}

				if tc.params.InTrash != nil {
					if *tc.params.InTrash {
						assert.NotNil(t, pld.DeletedAt, "DeletedAt should have been set")
					} else {
						assert.Nil(t, pld.DeletedAt, "DeletedAt should have been unset")
					}
				} else {
					assert.Nil(t, pld.DeletedAt, "DeletedAt should have not changed")
				}
			}
		})
	}
}

func TestIntegrationUpdateOrganization(t *testing.T) {
	dbCon := dependencies.DB
	defer lifecycle.PurgeModels(t, dbCon)

	newOrg := testorganizations.NewPersisted(t, dbCon, nil)
	orphan := testexperience.NewPersisted(t, dbCon, nil)
	orphan.Organization.DeletedAt = db.Now()
	orphan.Organization.Update(dbCon)

	params := &experience.UpdateParams{
		ID:             orphan.ID,
		OrganizationID: ptrs.NewString(newOrg.ID),
	}

	_, admSession := testauth.NewAdminAuth(t, dbCon)
	rec := callUpdate(t, params, httptests.NewRequestAuth(admSession))

	if assert.Equal(t, http.StatusOK, rec.Code) {
		var pld *experience.Payload
		if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
			t.Fatal(err)
		}

		// We make sure what should have not changed, hasn't
		assert.Equal(t, orphan.ID, pld.ID, "ID should have not changed")
		assert.Equal(t, orphan.JobTitle, pld.JobTitle, "JobTitle should have not changed")
		assert.Equal(t, orphan.Location, pld.Location, "Location should have not changed")
		assert.Equal(t, orphan.Description, pld.Description, "Description should have not changed")
		assert.Equal(t, orphan.StartDate.String(), pld.StartDate.String(), "StartDate should have not changed")
		assert.Nil(t, orphan.EndDate, "EndDate should have not changed")

		// We make sure the organization has changed
		assert.Equal(t, newOrg.ID, pld.Organization.ID, "OrganizationID does not match")
		assert.Equal(t, newOrg.ID, pld.Organization.ID, "Organization.ID does not match")
		assert.Equal(t, newOrg.Name, pld.Organization.Name, "Organization.ID does not match")
	}
}

func callUpdate(t *testing.T, params *experience.UpdateParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: experience.Endpoints[experience.EndpointUpdate],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}
