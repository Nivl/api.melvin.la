package experience_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/router/params"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteInvalidParams(t *testing.T) {
	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing ID",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "id",
			Sources: map[string]url.Values{
				"url": url.Values{},
				"form": url.Values{
					"current_password": []string{"password"},
				},
			},
		},
		{
			Description: "Should fail on invalid ID",
			MsgMatch:    params.ErrMsgInvalidUUID,
			FieldName:   "id",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"not-a-uuid"},
				},
				"form": url.Values{
					"current_password": []string{"password"},
				},
			},
		},
	}

	g := experience.Endpoints[experience.EndpointDelete].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestDeleteValidParams(t *testing.T) {
	testCases := []struct {
		description string
		sources     map[string]url.Values
	}{
		{
			"Should fail on blank password",
			map[string]url.Values{
				"url": url.Values{
					"id": []string{"48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := experience.Endpoints[experience.EndpointDelete]
			data, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)

			if data != nil {
				p := data.(*experience.DeleteParams)
				assert.Equal(t, tc.sources["url"].Get("id"), p.ID)
			}
		})
	}
}

func TestDeleteAccess(t *testing.T) {
	testCases := []testguard.AccessTestCase{
		{
			Description: "Should fail for anonymous users",
			User:        nil,
			ErrCode:     http.StatusUnauthorized,
		},
		{
			Description: "Should work for logged users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
			ErrCode:     http.StatusForbidden,
		},
		{
			Description: "Should work for logged users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0", IsAdmin: true},
			ErrCode:     0,
		},
	}

	g := experience.Endpoints[experience.EndpointDelete].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestDeleteHappyPath(t *testing.T) {
	handlerParams := &experience.DeleteParams{
		ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectDeletion()
	mockDB.ExpectGet("*experience.Experience", func(args mock.Arguments) {
		org := args.Get(0).(*experience.Experience)
		org.ID = "aa44ca86-553e-4e16-8c30-2e50e63f7eaa"
		org.JobTitle = "Google"
	})

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.On("NoContent").Return()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := experience.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestDeleteUnexistingExperience(t *testing.T) {
	handlerParams := &experience.DeleteParams{
		ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGetNotFound("*experience.Experience")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := experience.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	req.AssertExpectations(t)
	mockDB.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, httpErr.HTTPStatus())
}

func TestDeleteNoDBCon(t *testing.T) {
	handlerParams := &experience.DeleteParams{
		ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectDeletionError()
	mockDB.ExpectGet("*experience.Experience", func(args mock.Arguments) {
		org := args.Get(0).(*experience.Experience)
		org.ID = "aa44ca86-553e-4e16-8c30-2e50e63f7eaa"
		org.JobTitle = "Google"
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := experience.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}
