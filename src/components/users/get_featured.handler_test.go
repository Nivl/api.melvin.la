package users_test

import (
	"net/http"
	"testing"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/melvin-laplanche/ml-api/src/components/users/testusers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetFeaturedHappyPath(t *testing.T) {
	featuredUser := testusers.NewProfile()

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGet("*users.Profile", func(args mock.Arguments) {
		user := args.Get(0).(*users.Profile)
		*user = *featuredUser
	})

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectOk("*users.ProfilePayload", func(args mock.Arguments) {
		pld := args.Get(0).(*users.ProfilePayload)
		assert.Equal(t, featuredUser.User.ID, pld.ID, "The user ID should not have changed")
		assert.Equal(t, featuredUser.Name, pld.Name, "Name should not have changed")
		assert.Equal(t, *featuredUser.LinkedIn, pld.LinkedIn, "the LinkedIn id should not have changed")
		assert.Empty(t, pld.Email, "the email should not be returned to anyone")
	})

	// Mock the request & add expectations
	req := &mockrouter.HTTPRequest{}
	req.On("Response").Return(res)

	// call the handler
	err := users.GetFeatured(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestGetFeaturedUnexistingUser(t *testing.T) {
	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetNotFound("*users.Profile")

	// Mock the request & add expectations
	req := &mockrouter.HTTPRequest{}

	// call the handler
	err := users.GetFeatured(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, httpErr.HTTPStatus())
}
