package sessions_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/router/params"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/melvin-laplanche/ml-api/src/components/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteInvalidParams(t *testing.T) {
	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing token",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "token",
			Sources: map[string]url.Values{
				"url": url.Values{
					"token": []string{""},
				},
				"form": url.Values{},
			},
		},
		{
			Description: "Should fail on invalid token",
			MsgMatch:    params.ErrMsgInvalidUUID,
			FieldName:   "token",
			Sources: map[string]url.Values{
				"url": url.Values{
					"token": []string{"xxx-yyyy"},
				},
				"form": url.Values{},
			},
		},
	}

	g := sessions.Endpoints[sessions.EndpointDelete].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestDeleteValidParams(t *testing.T) {
	testCases := []struct {
		description string
		sources     map[string]url.Values
	}{
		{
			"Should work with only a valid ID",
			map[string]url.Values{
				"url": url.Values{
					"token": []string{"48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
				},
				"form": url.Values{},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := sessions.Endpoints[sessions.EndpointDelete]
			data, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)

			if data != nil {
				p := data.(*sessions.DeleteParams)
				assert.Equal(t, tc.sources["url"].Get("token"), p.Token)
			}
		})
	}
}

func TestDeleteAccess(t *testing.T) {
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

	g := sessions.Endpoints[sessions.EndpointDelete].Guard
	testguard.AccessTest(t, g, testCases)
}

// TestDeleteHappyPath test a user loging out (removing the current session)
func TestDeleteHappyPath(t *testing.T) {
	user := &auth.User{ID: "3e916798-a090-4f22-b1d1-04a63fbed6ef"}
	session := &auth.Session{ID: "3642e0e6-788e-4161-92dd-6c52ea823da9", UserID: user.ID}

	handlerParams := &sessions.DeleteParams{
		Token: session.ID,
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGet("*auth.Session", func(args mock.Arguments) {
		// return a session that match the session currently in use
		sess := args.Get(0).(*auth.Session)
		sess.ID = session.ID
		sess.UserID = session.UserID
	})
	// delete call
	mockDB.ExpectDeletion()

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.On("NoContent").Return()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(user)
	req.On("Session").Return(session)

	// call the handler
	err := sessions.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestDeleteOtherSession(t *testing.T) {
	handlerParams := &sessions.DeleteParams{
		Token:           "d1ba3e60-d674-47f0-9694-59fbda0fc659",
		CurrentPassword: "valid password",
	}

	// Generate a password for the user
	userPassword, err := auth.CryptPassword(handlerParams.CurrentPassword)
	assert.NoError(t, err)

	// defined the request data
	user := &auth.User{ID: "3e916798-a090-4f22-b1d1-04a63fbed6ef", Password: userPassword}
	session := &auth.Session{ID: "3642e0e6-788e-4161-92dd-6c52ea823da9", UserID: user.ID}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGet("*auth.Session", func(args mock.Arguments) {
		// returns a session that matched the params and that is attached to the current user
		sess := args.Get(0).(*auth.Session)
		sess.ID = handlerParams.Token
		sess.UserID = session.UserID
	})
	// delete call
	mockDB.ExpectDeletion()

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.On("NoContent").Return()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(user)
	req.On("Session").Return(session)

	// call the handler
	err = sessions.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestDeleteOtherSessionWrongPassword(t *testing.T) {
	handlerParams := &sessions.DeleteParams{
		Token:           "d1ba3e60-d674-47f0-9694-59fbda0fc659",
		CurrentPassword: "invalid password",
	}

	// Generate a password for the user
	userPassword, err := auth.CryptPassword("Valid password")
	assert.NoError(t, err)

	// defined the request data
	user := &auth.User{ID: "3e916798-a090-4f22-b1d1-04a63fbed6ef", Password: userPassword}
	session := &auth.Session{ID: "3642e0e6-788e-4161-92dd-6c52ea823da9", UserID: user.ID}

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(user)
	req.On("Session").Return(session)

	// call the handler
	err = sessions.Delete(req, &router.Dependencies{DB: nil})

	// Assert everything
	assert.Error(t, err, "the handler should not have fail")
	req.AssertExpectations(t)
	assert.Equal(t, http.StatusUnauthorized, httperr.Convert(err).Code(), "Should have fail with a 401")
}

func TestDeleteSomeonesSession(t *testing.T) {
	handlerParams := &sessions.DeleteParams{
		Token:           "d1ba3e60-d674-47f0-9694-59fbda0fc659",
		CurrentPassword: "valid password",
	}

	// Generate a password for the user
	userPassword, err := auth.CryptPassword(handlerParams.CurrentPassword)
	assert.NoError(t, err)

	// defined the request data
	user := &auth.User{ID: "3e916798-a090-4f22-b1d1-04a63fbed6ef", Password: userPassword}
	session := &auth.Session{ID: "3642e0e6-788e-4161-92dd-6c52ea823da9", UserID: user.ID}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGet("*auth.Session", func(args mock.Arguments) {
		// returns a session that matched the params and that is attached to an other user
		sess := args.Get(0).(*auth.Session)
		sess.ID = handlerParams.Token
		sess.UserID = "d15e8b30-69ad-405b-a0f0-0e298b994d89"
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(user)
	req.On("Session").Return(session)

	// call the handler
	err = sessions.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	assert.Equal(t, http.StatusNotFound, httperr.Convert(err).Code(), "Should have fail with a 404")
}

func TestDeleteUnexistingSession(t *testing.T) {
	handlerParams := &sessions.DeleteParams{
		Token:           "d1ba3e60-d674-47f0-9694-59fbda0fc659",
		CurrentPassword: "valid password",
	}

	// Generate a password for the user
	userPassword, err := auth.CryptPassword(handlerParams.CurrentPassword)
	assert.NoError(t, err)

	// defined the request data
	user := &auth.User{ID: "3e916798-a090-4f22-b1d1-04a63fbed6ef", Password: userPassword}
	session := &auth.Session{ID: "3642e0e6-788e-4161-92dd-6c52ea823da9", UserID: user.ID}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGetNotFound("*auth.Session")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(user)
	req.On("Session").Return(session)

	// call the handler
	err = sessions.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	assert.Equal(t, http.StatusNotFound, httperr.Convert(err).Code(), "Should have fail with a 404")
}
