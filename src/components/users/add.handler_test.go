package users_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/lib/pq"
	"github.com/melvin-laplanche/ml-api/src/components/users"
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
			"Should fail on missing name",
			"parameter missing: name",
			map[string]url.Values{
				"form": url.Values{
					"email":    []string{"email@valid.tld"},
					"password": []string{"password"},
				},
			},
		},
		{
			"Should fail on missing email",
			"parameter missing: email",
			map[string]url.Values{
				"form": url.Values{
					"name":     []string{"username"},
					"password": []string{"password"},
				},
			},
		},
		{
			"Should fail on missing password",
			"parameter missing: password",
			map[string]url.Values{
				"form": url.Values{
					"name":  []string{"username"},
					"email": []string{"email@valid.tld"},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := users.Endpoints[users.EndpointAdd]
			_, err := endpts.Guard.ParseParams(tc.sources)
			if assert.Error(t, err, "expected the guard to fail") {
				assert.True(t, strings.Contains(err.Error(), tc.msgMatch),
					"the error \"%s\" should contain the string \"%s\"", err.Error(), tc.msgMatch)
			}
		})
	}
}

func TestAddValidData(t *testing.T) {
	handlerParams := &users.AddParams{
		Name:     "username",
		Email:    "email@domain.tld",
		Password: "valid password",
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
	mockDB.Mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.On("Created", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
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
	err = users.Add(req, deps)

	// Assert everything
	assert.Nil(t, err, "the handler should not have fail")
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
	mockDB, err := mockdb.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.SQLDB.Close()
	deps := &router.Dependencies{
		DB: mockDB.DB,
	}
	conflictError := new(pq.Error)
	conflictError.Code = db.ErrDup
	conflictError.Message = "error: duplicate field"
	mockDB.Mock.ExpectExec("INSERT INTO").WillReturnError(conflictError)

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err = users.Add(req, deps)

	// Assert everything
	assert.Error(t, err)
	req.AssertExpectations(t)

	httpErr := httperr.Convert(err)
	assert.Equal(t, http.StatusConflict, httpErr.Code())
}
