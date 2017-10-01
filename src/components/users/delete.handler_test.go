package users_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-params"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/stretchr/testify/assert"
)

func TestDeleteInvalidParams(t *testing.T) {
	t.Parallel()

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
		{
			Description: "Should fail on missing password",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "current_password",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
				},
				"form": url.Values{},
			},
		},
		{
			Description: "Should fail on blank password",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "current_password",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
				},
				"form": url.Values{
					"current_password": []string{"      "},
				},
			},
		},
	}

	g := users.Endpoints[users.EndpointDelete].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestDeleteValidParams(t *testing.T) {
	t.Parallel()

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
				"form": url.Values{
					"current_password": []string{"password"},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := users.Endpoints[users.EndpointDelete]
			data, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)

			if data != nil {
				p := data.(*users.DeleteParams)
				assert.Equal(t, tc.sources["url"].Get("id"), p.ID)
				assert.Equal(t, tc.sources["form"].Get("current_password"), p.CurrentPassword)
			}

		})
	}
}

func TestDeleteAccess(t *testing.T) {
	t.Parallel()

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

	g := users.Endpoints[users.EndpointDelete].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestDeleteHappyPath(t *testing.T) {
	t.Parallel()

	handlerParams := &users.DeleteParams{
		ID:              "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		CurrentPassword: "valid password",
	}

	userPassword, err := auth.CryptPassword(handlerParams.CurrentPassword)
	assert.NoError(t, err)
	user := &auth.User{
		ID:       handlerParams.ID,
		Password: userPassword,
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectDeletion()

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.On("NoContent").Return()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(user)

	// call the handler
	err = users.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestDeleteInvalidPassword(t *testing.T) {
	t.Parallel()

	handlerParams := &users.DeleteParams{
		ID:              "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		CurrentPassword: "invalid password",
	}

	userPassword, err := auth.CryptPassword("valid password")
	assert.NoError(t, err)
	user := &auth.User{
		ID:       handlerParams.ID,
		Password: userPassword,
	}

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(user)

	// call the handler
	err = users.Delete(req, &router.Dependencies{DB: nil})

	// Assert everything
	assert.Error(t, err, "the handler should not have fail")
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusUnauthorized, httpErr.HTTPStatus())
}

func TestDeleteInvalidUser(t *testing.T) {
	t.Parallel()

	handlerParams := &users.DeleteParams{
		ID:              "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		CurrentPassword: "valid password",
	}

	userPassword, err := auth.CryptPassword("valid password")
	assert.NoError(t, err)
	user := &auth.User{
		// Different user ID
		ID:       "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9",
		Password: userPassword,
	}

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(user)

	// call the handler
	err = users.Delete(req, &router.Dependencies{DB: nil})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusForbidden, httpErr.HTTPStatus())
}

func TestDeleteNoDBConOnDelete(t *testing.T) {
	t.Parallel()

	handlerParams := &users.DeleteParams{
		ID:              "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		CurrentPassword: "valid password",
	}

	userPassword, err := auth.CryptPassword(handlerParams.CurrentPassword)
	assert.NoError(t, err)
	user := &auth.User{
		ID:       handlerParams.ID,
		Password: userPassword,
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectDeletionError()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(user)

	// call the handler
	err = users.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}
