package users_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func InvalidParams(t *testing.T) {
	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on no params",
			MsgMatch:    "parameter missing",
			Sources: map[string]url.Values{
				"form": url.Values{},
			},
		},
		{
			Description: "Should fail on missing name",
			MsgMatch:    "parameter missing: name",
			Sources: map[string]url.Values{
				"form": url.Values{
					"email":    []string{"email@valid.tld"},
					"password": []string{"password"},
				},
			},
		},
		{
			Description: "Should fail on missing email",
			MsgMatch:    "parameter missing: email",
			Sources: map[string]url.Values{
				"form": url.Values{
					"name":     []string{"username"},
					"password": []string{"password"},
				},
			},
		},
		{
			Description: "Should fail on missing password",
			MsgMatch:    "parameter missing: password",
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
	handlerParams := &users.AddParams{
		Name:     "username",
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectInsert("*auth.User")

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectCreated("*users.Payload", func(args mock.Arguments) {
		user := args.Get(0).(*users.Payload)
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
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestAddConflict(t *testing.T) {
	handlerParams := &users.AddParams{
		Name:     "username",
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectInsertConflict("*auth.User")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := httperr.Convert(err)
	assert.Equal(t, http.StatusConflict, httpErr.Code())
}
