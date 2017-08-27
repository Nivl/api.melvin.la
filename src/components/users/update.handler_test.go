package users_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/melvin-laplanche/ml-api/src/components/users/testusers"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/router/params"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-rest-tools/types/ptrs"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateInvalidParams(t *testing.T) {
	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing ID",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "id",
			Sources: map[string]url.Values{
				"url":  url.Values{},
				"form": url.Values{},
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
				"form": url.Values{},
			},
		},
	}

	g := users.Endpoints[users.EndpointUpdate].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestUpdateValidParams(t *testing.T) {
	testCases := []struct {
		description string
		sources     map[string]url.Values
	}{
		{
			"Should work with only a valid ID",
			map[string]url.Values{
				"url": url.Values{
					"id": []string{"48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
				},
				"form": url.Values{},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := users.Endpoints[users.EndpointUpdate]
			data, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)

			if data != nil {
				p := data.(*users.UpdateParams)
				assert.Equal(t, tc.sources["url"].Get("id"), p.ID)
			}
		})
	}
}

func TestUpdateAccess(t *testing.T) {
	testCases := []testguard.AccessTestCase{
		{
			Description: "Should fail for anonymous users",
			User:        nil,
			ErrCode:     http.StatusUnauthorized,
		},
		{
			Description: "Should work for logged users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
			ErrCode:     0,
		},
	}

	g := users.Endpoints[users.EndpointUpdate].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestUpdateHappyPath(t *testing.T) {
	profile := testusers.NewProfile()

	handlerParams := &users.UpdateParams{
		ID:               profile.User.ID,
		CurrentPassword:  "fake",
		Email:            "new_email@domain.tld",
		FacebookUsername: ptrs.NewString("new_username"),
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectUpdate("*auth.User")
	mockDB.ExpectUpdate("*users.Profile")
	mockDB.ExpectGet("*users.Profile", func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *profile
	})

	// Mock the response & add expectati ons
	res := new(mockrouter.HTTPResponse)
	res.ExpectOk("*users.ProfilePayload", func(args mock.Arguments) {
		data := args.Get(0).(*users.ProfilePayload)
		assert.Equal(t, profile.User.Name, data.Name, "the name should have not changed")
		assert.Equal(t, handlerParams.Email, data.Email, "email should have been updated")
		assert.Equal(t, *handlerParams.FacebookUsername, data.FacebookUsername, "FacebookUsername should have been updated")
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(profile.User)

	// call the handler
	err := users.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestUpdateInvalidPassword(t *testing.T) {
	profile := testusers.NewProfile()
	handlerParams := &users.UpdateParams{
		ID:              profile.UserID,
		CurrentPassword: "invalid password",
		NewPassword:     "new password",
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGet("*users.Profile", func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *profile
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(profile.User)

	// call the handler
	err := users.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should not have fail")
	req.AssertExpectations(t)
	mockDB.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusUnauthorized, httpErr.HTTPStatus())
}

func TestUpdateInvalidUser(t *testing.T) {
	handlerParams := &users.UpdateParams{
		ID:              "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		CurrentPassword: "valid password",
	}

	userPassword, err := auth.CryptPassword("valid password")
	assert.NoError(t, err)
	user := &auth.User{
		ID:       "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9",
		Password: userPassword,
	}

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(user)

	// call the handler
	err = users.Update(req, &router.Dependencies{})

	// Assert everything
	assert.Error(t, err, "the handler should not have fail")
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusForbidden, httpErr.HTTPStatus())
}
