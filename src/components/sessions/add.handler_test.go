package sessions_test

import (
	"database/sql"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/melvin-laplanche/ml-api/src/components/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddInvalidParams(t *testing.T) {
	testCases := []struct {
		description string
		msgMatch    string
		sources     map[string]url.Values
	}{
		{
			"Should fail on no params",
			"parameter missing",
			map[string]url.Values{
				"form": url.Values{},
			},
		},
		{
			"Should fail on missing email",
			"parameter missing: email",
			map[string]url.Values{
				"form": url.Values{
					"password": []string{"password"},
				},
			},
		},
		{
			"Should fail on missing password",
			"parameter missing: password",
			map[string]url.Values{
				"form": url.Values{
					"email": []string{"email@valid.tld"},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := sessions.Endpoints[sessions.EndpointAdd]
			_, err := endpts.Guard.ParseParams(tc.sources)
			assert.Error(t, err, "expected the guard to fail")
			assert.True(t, strings.Contains(err.Error(), tc.msgMatch),
				"the error \"%s\" should contain the string \"%s\"", err.Error(), tc.msgMatch)
		})
	}
}

func TestAddValidData(t *testing.T) {
	handlerParams := &sessions.AddParams{
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.On("NamedExec", mock.AnythingOfType("string"), mock.AnythingOfType("*auth.Session")).Return(nil, nil)
	getCall := mockDB.On("Get", mock.AnythingOfType("*auth.User"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
	getCall.Return(nil)
	getCall.Run(func(args mock.Arguments) {
		u := args.Get(0).(*auth.User)
		u.ID = "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9"

		var err error
		u.Password, err = auth.CryptPassword(handlerParams.Password)
		assert.NoError(t, err)
	})

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.On("Created", mock.AnythingOfType("*sessions.Payload")).Return(nil).Run(func(args mock.Arguments) {
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
	assert.Nil(t, err, "the handler should not have fail")
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
	mockDB := new(mockdb.DB)
	getCall := mockDB.On("Get", mock.AnythingOfType("*auth.User"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
	getCall.Return(sql.ErrNoRows)

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := sessions.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err)
	req.AssertExpectations(t)

	httpErr := httperr.Convert(err)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code())
}

func TestAddWrongPassword(t *testing.T) {
	handlerParams := &sessions.AddParams{
		Email:    "email@domain.tld",
		Password: "invalid password",
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	getCall := mockDB.On("Get", mock.AnythingOfType("*auth.User"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
	getCall.Return(nil)
	getCall.Run(func(args mock.Arguments) {
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

	httpErr := httperr.Convert(err)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code())
}
