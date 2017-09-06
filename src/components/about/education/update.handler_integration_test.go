// +build integration

package education_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dchest/uniuri"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/types/datetime"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/Nivl/go-rest-tools/types/ptrs"
	"github.com/melvin-laplanche/ml-api/src/components/about/education"
	"github.com/melvin-laplanche/ml-api/src/components/about/education/testeducation"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationUpdate(t *testing.T) {
	dbCon := deps.DB()

	defer lifecycle.PurgeModels(t, dbCon)
	_, admSession := testauth.NewAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)

	noop := testeducation.NewPersisted(t, dbCon, nil)
	changeAll := testeducation.NewPersisted(t, dbCon, nil)
	toUntrash := testeducation.NewPersisted(t, dbCon, &education.Education{
		DeletedAt: datetime.Now(),
	})

	tests := []struct {
		description string
		code        int
		toUpdate    *education.Education
		params      *education.UpdateParams
	}{
		{
			"Valid request should work",
			http.StatusOK,
			changeAll,
			&education.UpdateParams{
				ID:          changeAll.ID,
				Degree:      ptrs.NewString(uniuri.New()),
				GPA:         ptrs.NewString(uniuri.NewLen(4)),
				Location:    ptrs.NewString(uniuri.New()),
				Description: ptrs.NewString(uniuri.New()),
				StartYear:   ptrs.NewInt(2010),
				EndYear:     ptrs.NewInt(2013),
				InTrash:     ptrs.NewBool(true),
			},
		},
		{
			"Untrash should work",
			http.StatusOK,
			toUntrash,
			&education.UpdateParams{
				ID:      toUntrash.ID,
				InTrash: ptrs.NewBool(false),
			},
		},
		{
			"Noop should work",
			http.StatusOK,
			noop,
			&education.UpdateParams{
				ID: noop.ID,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callUpdate(t, tc.params, adminAuth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusOK {
				var pld *education.Payload
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tc.toUpdate.ID, pld.ID, "ID should have not changed")
				if tc.params.Degree != nil {
					assert.Equal(t, *tc.params.Degree, pld.Degree, "Degree should have changed")
				} else {
					assert.Equal(t, tc.toUpdate.Degree, pld.Degree, "Degree should have not changed")
				}

				if tc.params.GPA != nil {
					assert.Equal(t, *tc.params.GPA, pld.GPA, "GPA should have changed")
				} else {
					assert.Equal(t, tc.toUpdate.GPA, pld.GPA, "GPA should have not changed")
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

				if tc.params.StartYear != nil {
					assert.Equal(t, *tc.params.StartYear, pld.StartYear, "StartYear should have changed")
				} else {
					assert.Equal(t, tc.toUpdate.StartYear, pld.StartYear, "StartYear should have not changed")
				}

				if tc.params.EndYear != nil {
					assert.Equal(t, *tc.params.EndYear, pld.EndYear, "EndYear should have changed")
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
	dbCon := deps.DB()
	defer lifecycle.PurgeModels(t, dbCon)

	newOrg := testorganizations.NewPersisted(t, dbCon, nil)
	orphan := testeducation.NewPersisted(t, dbCon, nil)
	orphan.Organization.DeletedAt = datetime.Now()
	orphan.Organization.Update(dbCon)

	params := &education.UpdateParams{
		ID:             orphan.ID,
		OrganizationID: ptrs.NewString(newOrg.ID),
	}

	_, admSession := testauth.NewAdminAuth(t, dbCon)
	rec := callUpdate(t, params, httptests.NewRequestAuth(admSession))

	if assert.Equal(t, http.StatusOK, rec.Code) {
		var pld *education.Payload
		if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
			t.Fatal(err)
		}

		// We make sure what should have not changed, hasn't
		assert.Equal(t, orphan.ID, pld.ID, "ID should have not changed")
		assert.Equal(t, orphan.Degree, pld.Degree, "Degree should have not changed")
		assert.Equal(t, orphan.GPA, pld.GPA, "GPA should have not changed")
		assert.Equal(t, orphan.Location, pld.Location, "Location should have not changed")
		assert.Equal(t, orphan.Description, pld.Description, "Description should have not changed")
		assert.Equal(t, orphan.StartYear, pld.StartYear, "StartYear should have not changed")
		assert.Nil(t, orphan.EndYear, "EndYear should have not changed")

		// We make sure the organization has changed
		assert.Equal(t, newOrg.ID, pld.Organization.ID, "OrganizationID does not match")
		assert.Equal(t, newOrg.ID, pld.Organization.ID, "Organization.ID does not match")
		assert.Equal(t, newOrg.Name, pld.Organization.Name, "Organization.ID does not match")
	}
}

func callUpdate(t *testing.T, params *education.UpdateParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: education.Endpoints[education.EndpointUpdate],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}
