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

func TestIntegrationGetFeatured(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	featuredProfile, _ := testusers.NewAuthProfile(t, dbCon)
	featuredProfile.IsFeatured = ptrs.NewBool(true)
	assert.NoError(t, featuredProfile.Update(dbCon))

	// creates a bunch of users
	for i := 0; i < 5; i++ {
		testusers.NewPersistedProfile(t, dbCon, nil)
	}

	rec := callGetFeatured(t, helper.Deps)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var profile users.ProfilePayload
		if err := json.NewDecoder(rec.Body).Decode(&profile); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, featuredProfile.User.ID, profile.ID, "Wrong profile returned")
		assert.Empty(t, profile.Email, "the Email should be private")
		assert.Equal(t, *featuredProfile.LastName, profile.LastName, "Wrong LastName returned")
		assert.Equal(t, *featuredProfile.FirstName, profile.FirstName, "Wrong FirstName returned")
	}
}

func callGetFeatured(t *testing.T, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: users.Endpoints[users.EndpointGetFeatured],
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
