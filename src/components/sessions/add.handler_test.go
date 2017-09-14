package sessions_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-params"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/melvin-laplanche/ml-api/src/components/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInvalidParams(t *testing.T) {
	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing email",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "email",
			Sources: map[string]url.Values{
				"form": url.Values{
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
					"email": []string{"email@valid.tld"},
				},
			},
		},
	}

	g := sessions.Endpoints[sessions.EndpointAdd].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestAddValidData(t *testing.T) {
	handlerParams := &sessions.AddParams{
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectInsert("*auth.Session")
	mockDB.ExpectGet("*auth.User", func(args mock.Arguments) {
		u := args.Get(0).(*auth.User)
		u.ID = "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9"

		var err error
		u.Password, err = auth.CryptPassword(handlerParams.Password)
		assert.NoError(t, err)
	})

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectCreated("*sessions.Payload", func(args mock.Arguments) {
		session := args.Get(0).(*sessions.Payload)
		assert.Equal(t, "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9", session.UserID)
		assert.NotEmpty(t, session.Token)
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := sessions.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestAddUnexistingEmail(t *testing.T) {
	handlerParams := &sessions.AddParams{
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetNotFound("*auth.User")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := sessions.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusBadRequest, httpErr.HTTPStatus())
	assert.Equal(t, "email/password", httpErr.Field())
}

func TestAddWrongPassword(t *testing.T) {
	handlerParams := &sessions.AddParams{
		Email:    "email@domain.tld",
		Password: "invalid password",
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGet("*auth.User", func(args mock.Arguments) {
		u := args.Get(0).(*auth.User)
		u.ID = "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9"

		var err error
		u.Password, err = auth.CryptPassword("valid password")
		assert.NoError(t, err)
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := sessions.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusBadRequest, httpErr.HTTPStatus())
	assert.Equal(t, "email/password", httpErr.Field())
}

func TestAddNoDbConOnGet(t *testing.T) {
	handlerParams := &sessions.AddParams{
		Email:    "email@domain.tld",
		Password: "invalid password",
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetError("*auth.User")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := sessions.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}

func TestAddNoDBConOnSave(t *testing.T) {
	handlerParams := &sessions.AddParams{
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectInsertError("*auth.Session")
	mockDB.ExpectGet("*auth.User", func(args mock.Arguments) {
		u := args.Get(0).(*auth.User)
		u.ID = "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9"

		var err error
		u.Password, err = auth.CryptPassword(handlerParams.Password)
		assert.NoError(t, err)
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := sessions.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}
