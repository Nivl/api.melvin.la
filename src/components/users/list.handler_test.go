package users_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/melvin-laplanche/ml-api/src/components/users/testusers"

	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListInvalidParams(t *testing.T) {
	t.Parallel()

	testCases := []testguard.InvalidParamsTestCase{
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
	}

	g := users.Endpoints[users.EndpointList].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestListValidParams(t *testing.T) {
	t.Parallel()

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
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := users.Endpoints[users.EndpointList]
			_, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)
		})
	}
}

func TestListAccess(t *testing.T) {
	t.Parallel()

	testCases := []testguard.AccessTestCase{
		{
			Description: "Should fail for anonymous users",
			User:        nil,
			ErrCode:     http.StatusUnauthorized,
		},
		{
			Description: "Should fail for logged users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
			ErrCode:     http.StatusForbidden,
		},
		{
			Description: "Should work for admin users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0", IsAdmin: true},
			ErrCode:     0,
		},
	}

	g := users.Endpoints[users.EndpointList].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestListNoBDCon(t *testing.T) {
	t.Parallel()

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectSelectError("*users.Profiles")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(&users.ListParams{})

	// call the handler
	err := users.List(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}

func TestListInvalidSort(t *testing.T) {
	t.Parallel()

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectSelectError("*users.Profiles")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(&users.ListParams{Sort: "not_a_field"})

	// call the handler
	err := users.List(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusBadRequest, httpErr.HTTPStatus())
	assert.Equal(t, "sort", httpErr.Field())
}

func TestListPrivacy(t *testing.T) {
	t.Parallel()

	totalProfiles := 5

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectSelect("*users.Profiles", func(args mock.Arguments) {
		out := args.Get(0).(*users.Profiles)
		for i := 0; i < totalProfiles; i++ {
			*out = append(*out, testusers.NewProfile())
		}
	})

	res := &mockrouter.HTTPResponse{}
	res.ExpectOk("*users.ProfilesPayload", func(args mock.Arguments) {
		pld := args.Get(0).(*users.ProfilesPayload)
		require.Equal(t, totalProfiles, len(pld.Results), "Wrong number of profiles returned")
		for _, p := range pld.Results {
			assert.NotEmpty(t, p.Email, "The email should be visible")
		}
	})

	// Mock the request & add expectations
	req := &mockrouter.HTTPRequest{}
	req.On("Params").Return(&users.ListParams{})
	req.On("Response").Return(res)

	// call the handler
	err := users.List(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	res.AssertExpectations(t)
	req.AssertExpectations(t)
}
