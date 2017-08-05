// +build integration

package organizations_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/Nivl/go-rest-tools/router/formfile/testformfile"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/storage/fs"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationUploadHappyPath(t *testing.T) {
	cwd, _ := os.Getwd()
	dbCon := dependencies.DB

	defer lifecycle.PurgeModels(t, dbCon)
	_, admSession := testauth.NewAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)
	org := testorganizations.NewOrganization(t, dbCon, nil)

	tests := []struct {
		description string
		params      *organizations.UploadLogoParams
	}{
		{
			"Valid request should work",
			&organizations.UploadLogoParams{
				ID:   org.ID,
				Logo: testformfile.NewFormFile(t, cwd, "black_pixel.png"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			defer tc.params.Logo.File.Close()
			rec := callUploadLogo(t, tc.params, adminAuth)
			assert.Equal(t, http.StatusOK, rec.Code)

			if rec.Code == http.StatusOK {
				pld := &organizations.Payload{}
				if err := json.NewDecoder(rec.Body).Decode(pld); err != nil {
					assert.NoError(t, err)
					return
				}

				assert.Equal(t, org.ID, pld.ID, "ID should have not changed")
				assert.Equal(t, org.Name, pld.Name, "Name should have not changed")
				assert.NotEmpty(t, pld.Logo, "Logo should not be empty")
				assert.Equal(t, "/", string(pld.Logo[0]), "Logo should start by a \"/\" as the storage provider should have fallback to the FS one. Got %s", pld.Logo)
				exists, err := fs.FileExists(pld.Logo)
				assert.NoError(t, err)
				assert.True(t, exists, "The file should exist on the filesystem")
			}
		})
	}
}

func callUploadLogo(t *testing.T, params *organizations.UploadLogoParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: organizations.Endpoints[organizations.EndpointUploadLogo],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}
