package users_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetInvalidParams(t *testing.T) {
	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing ID",
			MsgMatch:    "parameter missing",
			FieldName:   "id",
			Sources: map[string]url.Values{
				"url": url.Values{},
			},
		},
		{
			Description: "Should fail on invalid ID",
			MsgMatch:    "not a valid uuid",
			FieldName:   "id",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"not-a-uuid"},
				},
			},
		},
	}

	g := users.Endpoints[users.EndpointGet].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestGetValidParams(t *testing.T) {
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

			endpts := users.Endpoints[users.EndpointGet]
			data, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)

			if data != nil {
				p := data.(*users.GetParams)
				assert.Equal(t, tc.sources["url"].Get("id"), p.ID)
			}
		})
	}
}

func TestGetOthersData(t *testing.T) {
	handlerParams := &users.GetParams{
		ID: "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9",
	}
	requester := &auth.User{
		ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
	}
	userToGet := &auth.User{
		ID:      handlerParams.ID,
		Name:    "user name",
		Email:   "email@domain.tld",
		IsAdmin: false,
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGet("*auth.User", func(args mock.Arguments) {
		user := args.Get(0).(*auth.User)
		user.ID = userToGet.ID
		user.Name = userToGet.Name
		user.Email = userToGet.Email
		user.IsAdmin = userToGet.IsAdmin
	})

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectOk("*users.Payload", func(args mock.Arguments) {
		pld := args.Get(0).(*users.Payload)
		assert.Equal(t, userToGet.ID, pld.ID, "ID should have not changed")
		assert.Equal(t, userToGet.Name, pld.Name, "Name should have not changed")
		assert.Empty(t, pld.Email, "the email should not be returned to anyone")
		assert.False(t, pld.IsAdmin, "user should not be an admin")
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(requester)

	// call the handler
	err := users.Get(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestGetOwnData(t *testing.T) {
	handlerParams := &users.GetParams{
		ID: "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9",
	}
	requester := &auth.User{
		ID:      handlerParams.ID,
		Name:    "user name",
		Email:   "email@domain.tld",
		IsAdmin: false,
	}

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectOk("*users.Payload", func(args mock.Arguments) {
		pld := args.Get(0).(*users.Payload)
		assert.Equal(t, requester.ID, pld.ID, "ID should have not changed")
		assert.Equal(t, requester.Name, pld.Name, "Name should have not changed")
		assert.Equal(t, requester.Email, pld.Email, "the email should be returned")
		assert.False(t, pld.IsAdmin, "user should not be an admin")
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(requester)

	// call the handler
	err := users.Get(req, &router.Dependencies{})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestGetUnexistingUser(t *testing.T) {
	handlerParams := &users.GetParams{
		ID: "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9",
	}
	requester := &auth.User{
		ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGetNotFound("*auth.User")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(requester)

	// call the handler
	err := users.Get(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := httperr.Convert(err)
	assert.Equal(t, http.StatusNotFound, httpErr.Code())
}
