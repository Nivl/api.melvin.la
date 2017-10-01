package users_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-params"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInvalidParams(t *testing.T) {
	t.Parallel()

	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing name",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "name",
			Sources: map[string]url.Values{
				"form": url.Values{
					"email":    []string{"email@valid.tld"},
					"password": []string{"password"},
				},
			},
		},
		{
			Description: "Should fail on missing email",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "email",
			Sources: map[string]url.Values{
				"form": url.Values{
					"name":     []string{"username"},
					"password": []string{"password"},
				},
			},
		},
		{
			Description: "Should fail on missing password",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "password",
			Sources: map[string]url.Values{
				"form": url.Values{
					"name":  []string{"username"},
					"email": []string{"email@valid.tld"},
				},
			},
		},
	}

	g := users.Endpoints[users.EndpointAdd].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestAddHappyPath(t *testing.T) {
	t.Parallel()

	handlerParams := &users.AddParams{
		Name:     "username",
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	tx, _ := mockDB.ExpectTransaction()
	tx.ExpectUpdate("*auth.User")
	tx.ExpectUpdate("*users.Profile")
	tx.ExpectCommit()
	tx.ExpectRollback()

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectCreated("*users.ProfilePayload", func(args mock.Arguments) {
		user := args.Get(0).(*users.ProfilePayload)
		assert.Equal(t, handlerParams.Name, user.Name)
		assert.Equal(t, handlerParams.Email, user.Email)
		assert.NotEmpty(t, user.ID)
		assert.False(t, user.IsAdmin)
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Nil(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	tx.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestAddConflict(t *testing.T) {
	t.Parallel()

	handlerParams := &users.AddParams{
		Name:     "username",
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	tx, _ := mockDB.ExpectTransaction()
	tx.ExpectUpdateConflict("*auth.User", "email")
	tx.ExpectRollback()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	tx.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusConflict, apiError.HTTPStatus())
	assert.Equal(t, "email", apiError.Field())
}

func TestAddProfileError(t *testing.T) {
	t.Parallel()

	handlerParams := &users.AddParams{
		Name:     "username",
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	tx, _ := mockDB.ExpectTransaction()
	tx.ExpectUpdate("*auth.User")
	tx.ExpectUpdateError("*users.Profile")
	tx.ExpectRollback()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	tx.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, apiError.HTTPStatus())
}

func TestAddCommitError(t *testing.T) {
	t.Parallel()

	handlerParams := &users.AddParams{
		Name:     "username",
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	tx, _ := mockDB.ExpectTransaction()
	tx.ExpectUpdate("*auth.User")
	tx.ExpectUpdate("*users.Profile")
	tx.ExpectCommitError()
	tx.ExpectRollback()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	tx.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, apiError.HTTPStatus())
}

func TestAddTransactionError(t *testing.T) {
	t.Parallel()

	handlerParams := &users.AddParams{
		Name:     "username",
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectTransactionError()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, apiError.HTTPStatus())
}
