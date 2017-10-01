// +build integration

package users_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/testing/integration"
	"github.com/Nivl/go-types/ptrs"
	"github.com/melvin-laplanche/ml-api/src/components/api"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/melvin-laplanche/ml-api/src/components/users/testusers"
	"github.com/stretchr/testify/assert"
)

func TestbatchUpdateIsFeatured(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	a1, as1 := testusers.NewAdminAuth(t, dbCon)
	params := &users.BatchUpdateParams{FeaturedUser: ptrs.NewString(a1.ID)}
	rec := callBatchUpdate(t, params, httptests.NewRequestAuth(as1), helper.Deps)

	if assert.Equal(t, http.StatusOK, rec.Code) {
		var pld users.ProfilesPayload
		if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
			t.Fatal(err)
		}

		if assert.Equal(t, 1, len(pld.Results), "Wrong number of results") {
			assert.True(t, pld.Results[0].IsFeatured, "The profile should be featured")
		}
	}
}

func TestbatchUpdateReplaceFeatured(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	a1, as1 := testusers.NewAdminAuth(t, dbCon)
	p1, _ := testusers.NewAuthProfile(t, dbCon)
	p1.IsFeatured = ptrs.NewBool(true)
	assert.NoError(t, p1.Update(dbCon), "Update failed")

	params := &users.BatchUpdateParams{FeaturedUser: ptrs.NewString(a1.ID)}
	rec := callBatchUpdate(t, params, httptests.NewRequestAuth(as1), helper.Deps)

	if assert.Equal(t, http.StatusOK, rec.Code) {
		var pld users.ProfilesPayload
		if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
			t.Fatal(err)
		}

		if assert.Equal(t, 2, len(pld.Results), "Wrong number of results") {
			currentFeatured := pld.Results[1]
			previousFeatured := pld.Results[0]

			if pld.Results[0].ID == a1.ID {
				currentFeatured = pld.Results[0]
				previousFeatured = pld.Results[1]
			}

			assert.True(t, currentFeatured.IsFeatured, "The profile should be featured")
			assert.Empty(t, previousFeatured.IsFeatured, "The profile should not be featured")
		}
	}
}

func callBatchUpdate(t *testing.T, params *users.BatchUpdateParams, auth *httptests.RequestAuth, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: users.Endpoints[users.EndpointBatchUpdate],
		Params:   params,
		Auth:     auth,
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
