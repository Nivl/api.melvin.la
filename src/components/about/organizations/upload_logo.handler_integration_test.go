// +build integration

package organizations_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Nivl/go-params/formfile/testformfile"
	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/storage/fs"
	"github.com/Nivl/go-rest-tools/testing/integration"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/melvin-laplanche/ml-api/src/components/api"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationUploadHappyPath(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	cwd, _ := os.Getwd()
	_, admSession := testauth.NewPersistedAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)
	org := testorganizations.NewPersisted(t, dbCon, nil)

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

	t.Run("parallel", func(t *testing.T) {
		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				defer helper.RecoverPanic()
				defer tc.params.Logo.File.Close()

				rec := callUploadLogo(t, tc.params, adminAuth, helper.Deps)
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
	})
}

func callUploadLogo(t *testing.T, params *organizations.UploadLogoParams, auth *httptests.RequestAuth, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: organizations.Endpoints[organizations.EndpointUploadLogo],
		Params:   params,
		Auth:     auth,
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
