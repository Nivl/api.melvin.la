// +build integration

package organizations_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/satori/go.uuid"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/testing/integration"
	"github.com/Nivl/go-types/datetime"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/melvin-laplanche/ml-api/src/components/api"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationListPagination tests the pagination
func TestIntegrationListPagination(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	_, admSession := testauth.NewPersistedAdminAuth(t, dbCon)
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

	t.Run("parallel", func(t *testing.T) {
		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				defer helper.RecoverPanic()

				rec := callList(t, tc.params, adminAuth, helper.Deps)
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
	})
}

// TestIntegrationListOrdering tests the ordering
func TestIntegrationListOrdering(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

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
	_, admSession := testauth.NewPersistedAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)

	// We set the default params manually otherwise it will send 0
	params := &organizations.ListParams{
		HandlerParams: paginator.HandlerParams{
			Page:    1,
			PerPage: 100,
		},
	}

	// make the request
	rec := callList(t, params, adminAuth, helper.Deps)

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

func callList(t *testing.T, params *organizations.ListParams, auth *httptests.RequestAuth, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: organizations.Endpoints[organizations.EndpointList],
		Params:   params,
		Auth:     auth,
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
