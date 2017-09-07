// +build integration

package organizations_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/satori/go.uuid"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/types/datetime"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationListPagination tests the pagination
func TestIntegrationListPagination(t *testing.T) {
	dbCon := deps.DB()

	defer lifecycle.PurgeModels(t, dbCon)
	_, admSession := testauth.NewAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)

	for i := 0; i < 35; i++ {
		testorganizations.NewPersisted(t, dbCon, nil)
	}
	// adding a deleted organization
	testorganizations.NewPersisted(t, dbCon, &organizations.Organization{
		DeletedAt: datetime.Now(),
	})

	tests := []struct {
		description   string
		expectedTotal int
		params        *organizations.ListParams
	}{
		{
			"100 per page with deleted",
			36,
			&organizations.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 100,
				},
				Deleted: true,
			},
		},
		{
			"100 per page without deleted",
			35,
			&organizations.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 100,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callList(t, tc.params, adminAuth)
			assert.Equal(t, http.StatusOK, rec.Code)

			if rec.Code == http.StatusOK {
				var pld organizations.ListPayload
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.expectedTotal, len(pld.Results), "invalid number of results")
			}
		})
	}
}

// TestIntegrationListOrdering tests the ordering
func TestIntegrationListOrdering(t *testing.T) {
	dbCon := deps.DB()
	defer lifecycle.PurgeModels(t, dbCon)

	// Creates the data
	names := []string{"z", "b", "y", "r", "a", "k", "f", "v"}
	// Add a uuid to the names so we avoid potential conflicts with other tests
	for i, _ := range names {
		names[i] += uuid.NewV4().String()
	}
	// create the orgs and save them to the database
	for _, name := range names {
		testorganizations.NewPersisted(t, dbCon, &organizations.Organization{
			Name: name,
		})
	}

	// the result should be sorted alphabetically
	expectedNames := make(sort.StringSlice, len(names))
	copy(expectedNames, names)
	expectedNames.Sort()

	// auth of the request
	_, admSession := testauth.NewAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)

	// We set the default params manually otherwise it will send 0
	params := &organizations.ListParams{
		HandlerParams: paginator.HandlerParams{
			Page:    1,
			PerPage: 100,
		},
	}

	// make the request
	rec := callList(t, params, adminAuth)

	// Assert everything went well
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var pld organizations.ListPayload
		if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
			t.Fatal(err)
		}

		// make sure we have the right number of results to avoid segfaults
		ok := assert.Equal(t, len(expectedNames), len(pld.Results), "invalid number of results")
		if ok {
			// assert the result has been ordered correctly
			for i, org := range pld.Results {
				assert.Equal(t, expectedNames[i], org.Name, "expected a different ordering")
			}
		}
	}
}

func callList(t *testing.T, params *organizations.ListParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: organizations.Endpoints[organizations.EndpointList],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}
