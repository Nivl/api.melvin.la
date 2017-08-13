// build integration

package experience_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/types/models/lifecycle"
	"github.com/Nivl/go-rest-tools/types/ptrs"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience/testexperience"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationListFiltering(t *testing.T) {
	dbCon := dependencies.DB
	defer lifecycle.PurgeModels(t, dbCon)

	totalBasicExp := 35
	for i := 0; i < totalBasicExp; i++ {
		testexperience.NewPersisted(t, dbCon, nil)
	}

	// adding a deleted experience
	testexperience.NewPersisted(t, dbCon, &experience.Experience{
		DeletedAt: db.Now(),
	})

	// adding an orphan experience
	orphan := testexperience.NewPersisted(t, dbCon, nil)
	orphan.Organization.DeletedAt = db.Now()
	orphan.Organization.Update(dbCon)

	// Adding an orphan that is also deleted
	orphanDeleted := testexperience.NewPersisted(t, dbCon, &experience.Experience{
		DeletedAt: db.Now(),
	})
	orphanDeleted.Organization.DeletedAt = db.Now()
	orphanDeleted.Organization.Update(dbCon)

	_, adminSession := testauth.NewAdminAuth(t, dbCon)

	tests := []struct {
		description   string
		expectedTotal int
		auth          *httptests.RequestAuth
		params        *experience.ListParams
	}{
		{
			"Admin default should returns everything",
			totalBasicExp + 3,
			httptests.NewRequestAuth(adminSession),
			&experience.ListParams{
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
			&experience.ListParams{
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
			&experience.ListParams{
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
			&experience.ListParams{
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
			&experience.ListParams{
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
			&experience.ListParams{
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
			&experience.ListParams{
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
			&experience.ListParams{
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
			&experience.ListParams{
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
				var pld experience.ListPayload
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.expectedTotal, len(pld.Results), "invalid number of results")
			}
		})
	}
}

func TestIntegrationListPagination(t *testing.T) {
	dbCon := dependencies.DB
	defer lifecycle.PurgeModels(t, dbCon)

	totalExp := 35
	for i := 0; i < totalExp; i++ {
		testexperience.NewPersisted(t, dbCon, nil)
	}

	tests := []struct {
		description   string
		expectedTotal int
		params        *experience.ListParams
	}{
		{
			"page 1, per_page 100",
			totalExp,
			&experience.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 100,
				},
			},
		},
		{
			"page 1, per_page 10",
			10,
			&experience.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    1,
					PerPage: 10,
				},
			},
		},
		{
			"page 2, per_page 34",
			1,
			&experience.ListParams{
				HandlerParams: paginator.HandlerParams{
					Page:    2,
					PerPage: 34,
				},
			},
		},
		{
			"page 4, per_page 10",
			5,
			&experience.ListParams{
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
				var pld experience.ListPayload
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.expectedTotal, len(pld.Results), "invalid number of results")
			}
		})
	}
}

func TestIntegrationListOrdering(t *testing.T) {
	dbCon := dependencies.DB
	defer lifecycle.PurgeModels(t, dbCon)

	dates := []struct {
		order int
		start string
		end   string
	}{
		{3, "2015-09", "2016-03"},
		{0, "2017-08", ""},
		{5, "2011-07", "2013-09"},
		{2, "2014-09", "2016-06"},
		{4, "2013-09", "2014-09"},
		{1, "2016-06", ""},
	}

	expWanted := make([]*experience.Experience, len(dates))

	for _, d := range dates {
		start, _ := db.NewDate(d.start)
		var end *db.Date
		if d.end != "" {
			end, _ = db.NewDate(d.end)
		}

		expWanted[d.order] = testexperience.NewPersisted(t, dbCon, &experience.Experience{
			StartDate: start,
			EndDate:   end,
		})
	}
	// We set the default params manually otherwise it will send 0
	params := &experience.ListParams{
		HandlerParams: paginator.HandlerParams{
			Page:    1,
			PerPage: 100,
		},
	}

	// make the request
	rec := callList(t, params, nil)

	// Assert everything went well
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var pld experience.ListPayload
		if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
			t.Fatal(err)
		}

		// make sure we have the right number of results to avoid segfaults
		ok := assert.Equal(t, len(expWanted), len(pld.Results), "invalid number of results")
		if ok {
			// assert the result has been ordered correctly
			for i, org := range pld.Results {
				assert.Equal(t, expWanted[i].ID, org.ID, "expected a different ordering")
			}
		}
	}
}

func callList(t *testing.T, params *experience.ListParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: experience.Endpoints[experience.EndpointList],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}
