// +build integration

package users_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/Nivl/go-rest-tools/types/ptrs"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/melvin-laplanche/ml-api/src/components/users/testusers"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationGetFeatured(t *testing.T) {
	dbCon := deps.DB
	defer lifecycle.PurgeModels(t, dbCon)

	featuredProfile, _ := testusers.NewAuthProfile(t, dbCon)
	featuredProfile.IsFeatured = ptrs.NewBool(true)
	assert.NoError(t, featuredProfile.Update(dbCon))

	// creates a bunch of users
	for i := 0; i < 5; i++ {
		testusers.NewPersistedProfile(t, dbCon, nil)
	}

	rec := callGetFeatured(t)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var profile users.ProfilePayload
		if err := json.NewDecoder(rec.Body).Decode(&profile); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, featuredProfile.User.ID, profile.ID, "Wrong profile returned")
		assert.Empty(t, profile.Email, "the Email should be private")
	}
}

func callGetFeatured(t *testing.T) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: users.Endpoints[users.EndpointGetFeatured],
	}

	return httptests.NewRequest(t, ri)
}
