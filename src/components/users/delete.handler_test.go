package users_test

import (
	"net/http"
	"net/url"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/stretchr/testify/assert"
)

func TestDeleteInvalidParams(t *testing.T) {
	testCases := []struct {
		description string
		sources     map[string]url.Values
	}{
		{
			"Should fail on no params",
			map[string]url.Values{
				"form": url.Values{},
			},
		},
		{
			"Should fail on missing ID",
			map[string]url.Values{
				"form": url.Values{
					"current_password": []string{"password"},
				},
			},
		},
		{
			"Should fail on invalid ID",
			map[string]url.Values{
				"form": url.Values{
					"id":               []string{"not-a-uuid"},
					"current_password": []string{"password"},
				},
			},
		},
		{
			"Should fail on missing password",
			map[string]url.Values{
				"form": url.Values{
					"id": []string{"48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
				},
			},
		},
		{
			"Should fail on blank password",
			map[string]url.Values{
				"form": url.Values{
					"id":               []string{"48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
					"current_password": []string{"      "},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := users.Endpoints[users.EndpointDelete]
			_, err := endpts.Guard.ParseParams(tc.sources)
			assert.Error(t, err)
		})
	}
}

func TestDeleteAccess(t *testing.T) {
	testCases := []struct {
		description string
		user        *auth.User
		errCode     int // <= 0 for no error
	}{
		{
			"Should fail for anonymous users",
			nil,
			http.StatusUnauthorized,
		},
		{
			"Should work for logged users",
			&auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
			0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := users.Endpoints[users.EndpointDelete]
			_, err := endpts.Guard.HasAccess(tc.user)
			if tc.errCode > 0 {
				assert.Error(t, err)
				assert.Equal(t, tc.errCode, httperr.Convert(err).Code())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteValidData(t *testing.T) {
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
	mockDB, err := mockdb.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.SQLDB.Close()
	deps := &router.Dependencies{
		DB: mockDB.DB,
	}
	mockDB.Mock.ExpectExec("DELETE").WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.On("NoContent").Return()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(user)

	// call the handler
	err = users.Delete(req, deps)

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestDeleteInvalidPassword(t *testing.T) {
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

	// Mock the database & add expectations
	deps := &router.Dependencies{
		DB: nil, // the DB shouldn't be used
	}

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(user)

	// call the handler
	err = users.Delete(req, deps)

	// Assert everything
	assert.Error(t, err, "the handler should not have fail")
	req.AssertExpectations(t)

	httpErr := httperr.Convert(err)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code())
}

func TestDeleteInvalidUser(t *testing.T) {
	handlerParams := &users.DeleteParams{
		ID:              "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		CurrentPassword: "valid password",
	}

	userPassword, err := auth.CryptPassword("valid password")
	assert.NoError(t, err)
	user := &auth.User{
		ID:       "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9",
		Password: userPassword,
	}

	// Mock the database & add expectations
	deps := &router.Dependencies{
		DB: nil, // the DB shouldn't be used
	}

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(user)

	// call the handler
	err = users.Delete(req, deps)

	// Assert everything
	assert.Error(t, err, "the handler should not have fail")
	req.AssertExpectations(t)

	httpErr := httperr.Convert(err)
	assert.Equal(t, http.StatusForbidden, httpErr.Code())
}
