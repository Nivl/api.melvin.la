// +build integration

package users_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/Nivl/go-rest-tools/security/auth"
	uuid "github.com/satori/go.uuid"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-types/datetime"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/melvin-laplanche/ml-api/src/components/users/testusers"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationListPagination tests the pagination
func TestIntegrationListPagination(t *testing.T) {
	dbCon := deps.DB()

	defer lifecycle.PurgeModels(t, dbCon)
	_, admSession := testauth.NewAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)

	for i := 0; i < 35; i++ {
		testusers.NewPersistedProfile(t, dbCon, nil)
	}
	// add a deleted user
	deleted := testusers.NewPersistedProfile(t, dbCon, nil)
	deleted.User.DeletedAt = datetime.Now()
	assert.NoError(t, deleted.User.Update(dbCon))

	tests := []struct {
		description   string
		expectedTotal int
		params        *users.ListParams
	}{
		{
			"100 per page",
			35,
			&users.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 100,
				},
			},
		},
		{
			"10 per page, page 1",
			10,
			&users.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 10,
				},
			},
		},
		{
			"10 per page, page 4",
			5,
			&users.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    4,
					PerPage: 10,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callList(t, tc.params, adminAuth)
			assert.Equal(t, http.StatusOK, rec.Code)

			if rec.Code == http.StatusOK {
				var pld users.ProfilesPayload
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.expectedTotal, len(pld.Results), "invalid number of results")
			}
		})
	}
}

// TestIntegrationListSorting tests the sorting
func TestIntegrationListSorting(t *testing.T) {
	dbCon := deps.DB()
	defer lifecycle.PurgeModels(t, dbCon)

	// Creates the data
	names := []string{"z", "b", "y", "r", "a", "k", "f", "v"}
	// Add a uuid to the names so we avoid potential conflicts with other tests
	for i, _ := range names {
		names[i] += uuid.NewV4().String()
	}
	// create the users and save them to the database
	for _, name := range names {
		testusers.NewPersistedProfile(t, dbCon, &users.Profile{
			User: &auth.User{
				Name: name,
			},
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
	params := &users.ListParams{
		HandlerParams: paginator.HandlerParams{
			Page:    1,
			PerPage: 100,
		},
		Sort: "name",
	}

	// make the request
	rec := callList(t, params, adminAuth)

	// Assert everything went well
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var pld users.ProfilesPayload
		if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
			t.Fatal(err)
		}

		// make sure we have the right number of results to avoid segfaults
		ok := assert.Equal(t, len(expectedNames), len(pld.Results), "invalid number of results")
		if ok {
			// assert the result has been ordered correctly
			for i, p := range pld.Results {
				assert.Equal(t, expectedNames[i], p.Name, "expected a different sorting")
			}
		}
	}
}

func callList(t *testing.T, params *users.ListParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: users.Endpoints[users.EndpointList],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}
