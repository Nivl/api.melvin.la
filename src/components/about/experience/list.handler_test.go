package experience_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/router/params"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience/testexperience"
	"github.com/stretchr/testify/assert"
)

func TestListInvalidParams(t *testing.T) {
	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on invalid delete value",
			MsgMatch:    params.ErrMsgInvalidBoolean,
			FieldName:   "deleted",
			Sources: map[string]url.Values{
				"query": url.Values{
					"deleted": []string{"not a bool"},
				},
			},
		},
		{
			Description: "Should fail on invalid orphans value",
			MsgMatch:    params.ErrMsgInvalidBoolean,
			FieldName:   "orphans",
			Sources: map[string]url.Values{
				"query": url.Values{
					"orphans": []string{"not a bool"},
				},
			},
		},
		{
			Description: "Should fail with page = 0",
			MsgMatch:    paginator.ErrMsgNumberBelow1,
			FieldName:   "page",
			Sources: map[string]url.Values{
				"query": url.Values{
					"page": []string{"0"},
				},
			},
		},
		{
			Description: "Should fail with per_page = 0",
			MsgMatch:    paginator.ErrMsgNumberBelow1,
			FieldName:   "per_page",
			Sources: map[string]url.Values{
				"query": url.Values{
					"per_page": []string{"0"},
				},
			},
		},
		{
			Description: "Should fail with per_page > 100",
			MsgMatch:    "cannot be > 100",
			FieldName:   "per_page",
			Sources: map[string]url.Values{
				"query": url.Values{
					"per_page": []string{"101"},
				},
			},
		},
		{
			Description: "Should fail with invalid operator",
			MsgMatch:    params.ErrMsgEnum,
			FieldName:   "op",
			Sources: map[string]url.Values{
				"query": url.Values{
					"op": []string{"nand"},
				},
			},
		},
	}

	g := experience.Endpoints[experience.EndpointList].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestListValidParams(t *testing.T) {
	testCases := []struct {
		description string
		sources     map[string]url.Values
	}{
		{
			"Should work with nothing",
			map[string]url.Values{
				"query": url.Values{},
			},
		},
		{
			"Should work with deleted=true",
			map[string]url.Values{
				"query": url.Values{
					"deleted": []string{"true"},
				},
			},
		},
		{
			"Should work with deleted=false",
			map[string]url.Values{
				"query": url.Values{
					"deleted": []string{"false"},
				},
			},
		},
		{
			"Should work with orphans=true",
			map[string]url.Values{
				"query": url.Values{
					"orphans": []string{"true"},
				},
			},
		},
		{
			"Should work with orphans=false",
			map[string]url.Values{
				"query": url.Values{
					"orphans": []string{"false"},
				},
			},
		},
		{
			"Should work with both orphans and deleted",
			map[string]url.Values{
				"query": url.Values{
					"orphans": []string{"true"},
					"deleted": []string{"true"},
				},
			},
		},
		{
			"Should work with orphans, deleted and an op",
			map[string]url.Values{
				"query": url.Values{
					"orphans": []string{"true"},
					"deleted": []string{"true"},
					"op":      []string{"or"},
				},
			},
		},
		{
			"Should work with a page",
			map[string]url.Values{
				"query": url.Values{
					"page": []string{"1"},
				},
			},
		},
		{
			"Should work with a per_page",
			map[string]url.Values{
				"query": url.Values{
					"per_page": []string{"10"},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := experience.Endpoints[experience.EndpointList]
			_, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)
		})
	}
}

func TestListAccess(t *testing.T) {
	testCases := []testguard.AccessTestCase{
		{
			Description: "Should work for anonymous users",
			User:        nil,
			ErrCode:     0,
		},
		{
			Description: "Should work for logged users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
			ErrCode:     0,
		},
		{
			Description: "Should work for admin users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0", IsAdmin: true},
			ErrCode:     0,
		},
	}

	g := experience.Endpoints[experience.EndpointList].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestListNoBDCon(t *testing.T) {
	requester := testauth.NewUser()

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectSelectError("*experience.ListExperience")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(&experience.ListParams{})
	req.On("User").Return(requester)

	// call the handler
	err := experience.List(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}

func TestListPublic(t *testing.T) {
	requester := testauth.NewUser()

	expList := experience.ListExperience{
		testexperience.New(),
		testexperience.New(),
		testexperience.New(),
		testexperience.New(),
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectSelect("*experience.ListExperience", func(args mock.Arguments) {
		l := args.Get(0).(*experience.ListExperience)
		*l = expList
	})

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectOk("*experience.ListPayload", func(args mock.Arguments) {
		l := args.Get(0).(*experience.ListPayload)
		if assert.Equal(t, len(expList), len(l.Results), "different number of results") {
			// make sure the results are public
			for i, r := range l.Results {
				assert.Equal(t, expList[i].ID, r.ID, "ID should have not changed")
				assert.Equal(t, expList[i].JobTitle, r.JobTitle, "JobTitle should have not changed")
				assert.Equal(t, expList[i].Location, r.Location, "Location should have not changed")
				assert.Equal(t, expList[i].Description, r.Description, "Description should have not changed")

				assert.Nil(t, r.CreatedAt, "CreatedAt should have not changed")
				assert.Nil(t, r.UpdatedAt, "UpdatedAt should have not changed")
			}
		}
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(&experience.ListParams{})
	req.On("User").Return(requester)
	req.On("Response").Return(res)

	// call the handler
	err := experience.List(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestListPrivate(t *testing.T) {
	requester := testauth.NewAdmin()

	expList := experience.ListExperience{
		testexperience.New(),
		testexperience.New(),
		testexperience.New(),
		testexperience.New(),
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectSelect("*experience.ListExperience", func(args mock.Arguments) {
		l := args.Get(0).(*experience.ListExperience)
		*l = expList
	})

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectOk("*experience.ListPayload", func(args mock.Arguments) {
		l := args.Get(0).(*experience.ListPayload)
		if assert.Equal(t, len(expList), len(l.Results), "different number of results") {
			// make sure the results are public
			for i, r := range l.Results {
				assert.Equal(t, expList[i].ID, r.ID, "ID should have not changed")
				assert.Equal(t, expList[i].JobTitle, r.JobTitle, "JobTitle should have not changed")
				assert.Equal(t, expList[i].Location, r.Location, "Location should have not changed")
				assert.Equal(t, expList[i].Description, r.Description, "Description should have not changed")

				assert.NotNil(t, r.CreatedAt, "CreatedAt should have not changed")
				assert.NotNil(t, r.UpdatedAt, "UpdatedAt should have not changed")
			}
		}
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(&experience.ListParams{})
	req.On("User").Return(requester)
	req.On("Response").Return(res)

	// call the handler
	err := experience.List(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}
