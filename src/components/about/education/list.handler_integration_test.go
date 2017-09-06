// +build integration

package education_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/types/datetime"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/Nivl/go-rest-tools/types/ptrs"
	"github.com/melvin-laplanche/ml-api/src/components/about/education"
	"github.com/melvin-laplanche/ml-api/src/components/about/education/testeducation"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationListFiltering(t *testing.T) {
	dbCon := deps.DB()
	defer lifecycle.PurgeModels(t, dbCon)

	totalBasicExp := 35
	for i := 0; i < totalBasicExp; i++ {
		testeducation.NewPersisted(t, dbCon, nil)
	}

	// adding a deleted education
	testeducation.NewPersisted(t, dbCon, &education.Education{
		DeletedAt: datetime.Now(),
	})

	// adding an orphan education
	orphan := testeducation.NewPersisted(t, dbCon, nil)
	orphan.Organization.DeletedAt = datetime.Now()
	orphan.Organization.Update(dbCon)

	// Adding an orphan that is also deleted
	orphanDeleted := testeducation.NewPersisted(t, dbCon, &education.Education{
		DeletedAt: datetime.Now(),
	})
	orphanDeleted.Organization.DeletedAt = datetime.Now()
	orphanDeleted.Organization.Update(dbCon)

	_, adminSession := testauth.NewAdminAuth(t, dbCon)

	tests := []struct {
		description   string
		expectedTotal int
		auth          *httptests.RequestAuth
		params        *education.ListParams
	}{
		{
			"Admin default should returns everything",
			totalBasicExp + 3,
			httptests.NewRequestAuth(adminSession),
			&education.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 100,
				},
			},
		},
		{
			"Admin: no deleted no orphans",
			totalBasicExp,
			httptests.NewRequestAuth(adminSession),
			&education.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 100,
				},
				Orphans: ptrs.NewBool(false),
				Deleted: ptrs.NewBool(false),
			},
		},
		{
			"Anonymous default should not return deleted and orphans",
			totalBasicExp,
			nil,
			&education.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 100,
				},
			},
		},
		{
			"Anonymous should not be able to filter deleted and orphans",
			totalBasicExp,
			nil,
			&education.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 100,
				},
				Orphans:  ptrs.NewBool(true),
				Deleted:  ptrs.NewBool(true),
				Operator: "or",
			},
		},
		{
			"Admin: all orphans",
			2,
			httptests.NewRequestAuth(adminSession),
			&education.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 100,
				},
				Orphans: ptrs.NewBool(true),
			},
		},
		{
			"Admin: all orphans BUT not deleted",
			1,
			httptests.NewRequestAuth(adminSession),
			&education.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 100,
				},
				Orphans: ptrs.NewBool(true),
				Deleted: ptrs.NewBool(false),
			},
		},
		{
			"Admin: All deleted",
			2,
			httptests.NewRequestAuth(adminSession),
			&education.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 100,
				},
				Deleted: ptrs.NewBool(true),
			},
		},
		{
			"Admin: all orphans AND deleted",
			1,
			httptests.NewRequestAuth(adminSession),
			&education.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 100,
				},
				Orphans: ptrs.NewBool(true),
				Deleted: ptrs.NewBool(true),
			},
		},
		{
			"Admin: all orphans OR deleted",
			3,
			httptests.NewRequestAuth(adminSession),
			&education.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 100,
				},
				Orphans:  ptrs.NewBool(true),
				Deleted:  ptrs.NewBool(true),
				Operator: "or",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callList(t, tc.params, tc.auth)
			assert.Equal(t, http.StatusOK, rec.Code)

			if rec.Code == http.StatusOK {
				var pld education.ListPayload
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.expectedTotal, len(pld.Results), "invalid number of results")
			}
		})
	}
}

func TestIntegrationListPagination(t *testing.T) {
	dbCon := deps.DB()
	defer lifecycle.PurgeModels(t, dbCon)

	totalExp := 35
	for i := 0; i < totalExp; i++ {
		testeducation.NewPersisted(t, dbCon, nil)
	}

	tests := []struct {
		description   string
		expectedTotal int
		params        *education.ListParams
	}{
		{
			"page 1, per_page 100",
			totalExp,
			&education.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 100,
				},
			},
		},
		{
			"page 1, per_page 10",
			10,
			&education.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 10,
				},
			},
		},
		{
			"page 2, per_page 34",
			1,
			&education.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    2,
					PerPage: 34,
				},
			},
		},
		{
			"page 4, per_page 10",
			5,
			&education.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    4,
					PerPage: 10,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callList(t, tc.params, nil)
			assert.Equal(t, http.StatusOK, rec.Code)

			if rec.Code == http.StatusOK {
				var pld education.ListPayload
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.expectedTotal, len(pld.Results), "invalid number of results")
			}
		})
	}
}

func TestIntegrationListOrdering(t *testing.T) {
	dbCon := deps.DB()
	defer lifecycle.PurgeModels(t, dbCon)

	dates := []struct {
		order int
		start int
		end   int
	}{
		{2, 2015, 2016},
		{0, 2017, -1},
		{5, 2011, 2013},
		{3, 2014, 2016},
		{4, 2013, 2014},
		{1, 2016, -1},
	}

	eduWanted := make([]*education.Education, len(dates))
	for _, d := range dates {
		var end *int
		if d.end != -1 {
			end = ptrs.NewInt(d.end)
		}

		eduWanted[d.order] = testeducation.NewPersisted(t, dbCon, &education.Education{
			StartYear: d.start,
			EndYear:   end,
		})
	}

	// We set the default params manually otherwise it will send 0
	params := &education.ListParams{
		HandlerParams: paginator.HandlerParams{
			Page:    1,
			PerPage: 100,
		},
	}

	// make the request
	rec := callList(t, params, nil)

	// Assert everything went well
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var pld education.ListPayload
		if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
			t.Fatal(err)
		}

		// make sure we have the right number of results to avoid segfaults
		ok := assert.Equal(t, len(eduWanted), len(pld.Results), "invalid number of results")
		if ok {
			// assert the result has been ordered correctly
			for i, org := range pld.Results {
				assert.Equal(t, eduWanted[i].ID, org.ID, "expected a different ordering")
			}
		}
	}
}

func callList(t *testing.T, params *education.ListParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: education.Endpoints[education.EndpointList],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}
