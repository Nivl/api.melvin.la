package users

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/stretchr/testify/assert"
)

func TestListInvalidParams(t *testing.T) {
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

	g := Endpoints[EndpointList].Guard
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
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := Endpoints[EndpointList]
			_, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)
		})
	}
}

func TestListAccess(t *testing.T) {
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

	g := Endpoints[EndpointList].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestListParamsGetSort(t *testing.T) {
	// sugar
	shouldFail := true

	testCases := []struct {
		description string
		fields      string
		expected    string
		shouldFail  bool
	}{
		{
			"No fields should return the default sorting",
			"",
			"is_featured ASC,created_at ASC",
			!shouldFail,
		},
		{
			"Order by ,, should return the default sorting",
			",,",
			"is_featured ASC,created_at ASC",
			!shouldFail,
		},
		{
			"Order by ,,,,,,, should return the default sorting",
			",,,,,,,",
			"is_featured ASC,created_at ASC",
			!shouldFail,
		},
		{
			"Order by ,,,name,,,, should sort by name",
			",,,name,,,,",
			"name ASC",
			!shouldFail,
		},
		{
			"Order by -name should work",
			"-name",
			"name DESC",
			!shouldFail,
		},
		{
			"Order by is_featured and -name should work",
			"is_featured,-name",
			"is_featured ASC,name DESC",
			!shouldFail,
		},
		{
			"Order by not_a_field should fail",
			"not_a_field",
			"",
			shouldFail,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			output, err := listParamsGetSort(tc.fields)
			if tc.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, output)
			}
		})
	}
}

func TestListNoBDCon(t *testing.T) {
	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectSelectError("*users.Profiles")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(&ListParams{})

	// call the handler
	err := List(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}

func TestListInvalidSort(t *testing.T) {
	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectSelectError("*users.Profiles")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(&ListParams{Sort: "not_a_field"})

	// call the handler
	err := List(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusBadRequest, httpErr.HTTPStatus())
	assert.Equal(t, "sort", httpErr.Field())
}
