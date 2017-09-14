// +build integration

package organizations_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-types/datetime"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/Nivl/go-types/ptrs"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationUpdate(t *testing.T) {
	dbCon := deps.DB()

	defer lifecycle.PurgeModels(t, dbCon)
	_, admSession := testauth.NewAdminAuth(t, dbCon)
	adminAuth := httptests.NewRequestAuth(admSession)

	noop := testorganizations.NewPersisted(t, dbCon, nil)
	changeAll := testorganizations.NewPersisted(t, dbCon, nil)
	toUntrash := testorganizations.NewPersisted(t, dbCon, &organizations.Organization{
		DeletedAt: datetime.Now(),
	})

	tests := []struct {
		description string
		code        int
		toUpdate    *organizations.Organization
		params      *organizations.UpdateParams
	}{
		{
			"Valid request should work",
			http.StatusOK,
			changeAll,
			&organizations.UpdateParams{
				ID:        changeAll.ID,
				Name:      ptrs.NewString("new name"),
				ShortName: ptrs.NewString("new short name"),
				Website:   ptrs.NewString("http://google.com"),
				InTrash:   ptrs.NewBool(true),
			},
		},
		{
			"Untrash should work",
			http.StatusOK,
			toUntrash,
			&organizations.UpdateParams{
				ID:        toUntrash.ID,
				InTrash:   ptrs.NewBool(false),
				Name:      nil,
				ShortName: nil,
				Website:   nil,
			},
		},
		{
			"Noop should work",
			http.StatusOK,
			noop,
			&organizations.UpdateParams{
				ID:        noop.ID,
				Name:      nil,
				ShortName: nil,
				Website:   nil,
				InTrash:   nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callUpdate(t, tc.params, adminAuth)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusOK {
				var pld *organizations.Payload
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tc.toUpdate.ID, pld.ID, "ID should have not changed")
				if tc.params.Name != nil {
					assert.Equal(t, *tc.params.Name, pld.Name, "Name should have changed")
				} else {
					assert.Equal(t, tc.toUpdate.Name, pld.Name, "Name should have not changed")
				}

				if tc.params.ShortName != nil {
					assert.Equal(t, *tc.params.ShortName, pld.ShortName, "ShortName should have changed")
				} else {
					assert.Equal(t, *tc.toUpdate.ShortName, pld.ShortName, "ShortName should have not changed")
				}

				if tc.params.Website != nil {
					assert.Equal(t, *tc.params.Website, pld.Website, "Website should have changed")
				} else {
					assert.Equal(t, *tc.toUpdate.Website, pld.Website, "Website should have not changed")
				}

				if tc.params.InTrash != nil {
					if *tc.params.InTrash {
						assert.NotNil(t, pld.DeletedAt, "DeletedAt should have been set")
					} else {
						assert.Nil(t, pld.DeletedAt, "DeletedAt should have been unset")
					}
				} else {
					assert.Nil(t, pld.DeletedAt, "DeletedAt should have not changed")
				}
			}
		})
	}
}

func callUpdate(t *testing.T, params *organizations.UpdateParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: organizations.Endpoints[organizations.EndpointUpdate],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}
