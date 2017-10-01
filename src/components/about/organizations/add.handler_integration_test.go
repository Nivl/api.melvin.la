// +build integration

package organizations_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/testing/integration"
	"github.com/Nivl/go-types/ptrs"
	"github.com/dchest/uniuri"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/melvin-laplanche/ml-api/src/components/api"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	_, admSession := testauth.NewPersistedAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)

	tests := []struct {
		description string
		code        int
		params      *organizations.AddParams
	}{
		{
			"Valid Request should work",
			http.StatusCreated,
			&organizations.AddParams{Name: uniuri.New(), Website: ptrs.NewString("http://google.com")},
		},
	}

	t.Run("parallel", func(t *testing.T) {
		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				defer helper.RecoverPanic()

				rec := callAdd(t, tc.params, adminAuth, helper.Deps)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusCreated {
					org := &organizations.Payload{}
					if err := json.NewDecoder(rec.Body).Decode(org); err != nil {
						t.Fatal(err)
					}

					assert.NotEmpty(t, org.ID)
					assert.Equal(t, tc.params.Name, org.Name)
					assert.Empty(t, org.ShortName)
					assert.Equal(t, *tc.params.Website, org.Website)
					assert.NotNil(t, org.CreatedAt)
				}
			})
		}
	})
}

func TestIntegrationAddConflictName(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	_, admSession := testauth.NewPersistedAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)
	org := testorganizations.NewPersisted(t, dbCon, nil)

	p := &organizations.AddParams{
		Name: org.Name,
	}

	rec := callAdd(t, p, adminAuth, helper.Deps)
	assert.Equal(t, http.StatusConflict, rec.Code)

	pld := &router.ResponseError{}
	if err := json.NewDecoder(rec.Body).Decode(pld); err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, org.ID)
	assert.Equal(t, "already exists", pld.Error)
	assert.Equal(t, "name", pld.Field)
}

func TestIntegrationAddConflictShortName(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	_, admSession := testauth.NewPersistedAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)
	org := testorganizations.NewPersisted(t, dbCon, nil)

	p := &organizations.AddParams{
		Name:      uniuri.New(),
		ShortName: org.ShortName,
	}

	rec := callAdd(t, p, adminAuth, helper.Deps)
	assert.Equal(t, http.StatusConflict, rec.Code)

	pld := &router.ResponseError{}
	if err := json.NewDecoder(rec.Body).Decode(pld); err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, org.ID)
	assert.Equal(t, "already exists", pld.Error)
	assert.Equal(t, "short_name", pld.Field)
}

func callAdd(t *testing.T, params *organizations.AddParams, auth *httptests.RequestAuth, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: organizations.Endpoints[organizations.EndpointAdd],
		Params:   params,
		Auth:     auth,
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
