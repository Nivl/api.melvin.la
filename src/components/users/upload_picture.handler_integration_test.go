// +build integration

package users_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/router/formfile/testformfile"
	"github.com/Nivl/go-rest-tools/storage/fs"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/melvin-laplanche/ml-api/src/components/users/testusers"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationUploadHappyPath(t *testing.T) {
	cwd, _ := os.Getwd()
	dbCon := dependencies.DB

	defer lifecycle.PurgeModels(t, dbCon)
	_, admSession := testusers.NewAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)
	userProfile := testusers.NewPersistedProfile(t, dbCon, nil)

	tests := []struct {
		description string
		params      *users.UploadPictureParams
	}{
		{
			"Valid request should work",
			&users.UploadPictureParams{
				ID:      userProfile.User.ID,
				Picture: testformfile.NewFormFile(t, cwd, "black_pixel.png"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			defer tc.params.Picture.File.Close()
			rec := callUploadPicture(t, tc.params, adminAuth)
			assert.Equal(t, http.StatusOK, rec.Code)

			if rec.Code == http.StatusOK {
				pld := &users.ProfilePayload{}
				if err := json.NewDecoder(rec.Body).Decode(pld); err != nil {
					assert.NoError(t, err)
					return
				}

				assert.Equal(t, userProfile.User.ID, pld.ID, "ID should have not changed")
				assert.NotEmpty(t, userProfile.Picture, "Picture should not be empty")
				assert.Equal(t, "/", string(pld.Picture[0]), "Picture should start by a \"/\" as the storage provider should have fallback to the FS one. Got %s", pld.Picture)
				exists, err := fs.FileExists(pld.Picture)
				assert.NoError(t, err)
				assert.True(t, exists, "The file should exist on the filesystem")
			}
		})
	}
}

func callUploadPicture(t *testing.T, params *users.UploadPictureParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: users.Endpoints[users.EndpointUploadPicture],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}
