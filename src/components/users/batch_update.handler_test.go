package users_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/melvin-laplanche/ml-api/src/components/users/testusers"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-params"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-types/ptrs"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBatchUpdateInvalidParams(t *testing.T) {
	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on invalid featured_user",
			MsgMatch:    params.ErrMsgInvalidUUID,
			FieldName:   "featured_user",
			Sources: map[string]url.Values{
				"form": url.Values{
					"featured_user": []string{"not-a-uuid"},
				},
			},
		},
		{
			Description: "Should fail on empty featured_user",
			MsgMatch:    params.ErrMsgEmptyParameter,
			FieldName:   "featured_user",
			Sources: map[string]url.Values{
				"form": url.Values{
					"featured_user": []string{""},
				},
			},
		},
	}

	g := users.Endpoints[users.EndpointBatchUpdate].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestBatchUpdateValidParams(t *testing.T) {
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

func TestBatchUpdateAccess(t *testing.T) {
	testCases := []testguard.AccessTestCase{
		{
			Description: "Should fail for anonymous users",
			User:        nil,
			ErrCode:     http.StatusUnauthorized,
		},
		{
			Description: "Should work for logged users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
			ErrCode:     http.StatusForbidden,
		},
		{
			Description: "Should work for admins",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0", IsAdmin: true},
			ErrCode:     0,
		},
	}

	g := users.Endpoints[users.EndpointBatchUpdate].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestBatchUpdateFirstFeaturedUser(t *testing.T) {
	profileToFeature := testusers.NewProfile()

	handlerParams := &users.BatchUpdateParams{
		FeaturedUser: ptrs.NewString(profileToFeature.User.ID),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetID("*users.Profile", profileToFeature.User.ID, func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *profileToFeature
	})
	mockDB.ExpectGetNoParamsNotFound("*users.Profile")
	mockDB.ExpectUpdate("*users.Profile")

	// Mock the response & add expectati ons
	res := new(mockrouter.HTTPResponse)
	res.ExpectOk("*users.ProfilesPayload", func(args mock.Arguments) {
		data := args.Get(0).(*users.ProfilesPayload)
		if assert.Equal(t, 1, len(data.Results), "wrong number of results") {
			assert.Equal(t, profileToFeature.User.ID, data.Results[0].ID, "Wrong profile returned")
			assert.True(t, data.Results[0].IsFeatured, "IsFeatured should be set to true")

			// check privacy
			assert.Equal(t, profileToFeature.User.Email, data.Results[0].Email, "Email should contain the user email")
		}
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.BatchUpdate(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestBatchUpdateFirstFeaturedUserSaveFail(t *testing.T) {
	profileToFeature := testusers.NewProfile()

	handlerParams := &users.BatchUpdateParams{
		FeaturedUser: ptrs.NewString(profileToFeature.User.ID),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetID("*users.Profile", profileToFeature.User.ID, func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *profileToFeature
	})
	mockDB.ExpectGetNoParamsNotFound("*users.Profile")
	mockDB.ExpectUpdateError("*users.Profile")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.BatchUpdate(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}

func TestBatchUpdateFeaturedUser(t *testing.T) {
	profileToFeature := testusers.NewProfile()

	currentFeatured := testusers.NewProfile()
	currentFeatured.IsFeatured = ptrs.NewBool(true)

	handlerParams := &users.BatchUpdateParams{
		FeaturedUser: ptrs.NewString(profileToFeature.User.ID),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetID("*users.Profile", profileToFeature.User.ID, func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *profileToFeature
	})
	mockDB.ExpectGetNoParams("*users.Profile", func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *currentFeatured
	})
	mockTX, _ := mockDB.ExpectTransaction()
	mockTX.ExpectUpdate("*users.Profile")
	mockTX.ExpectUpdate("*users.Profile")
	mockTX.ExpectCommit()
	mockTX.ExpectRollback()

	// Mock the response & add expectati ons
	res := new(mockrouter.HTTPResponse)
	res.ExpectOk("*users.ProfilesPayload", func(args mock.Arguments) {
		data := args.Get(0).(*users.ProfilesPayload)
		if assert.Equal(t, 2, len(data.Results), "wrong number of results") {
			if data.Results[0].ID == currentFeatured.User.ID {
				assert.False(t, data.Results[0].IsFeatured, "IsFeatured should not be set")
				assert.Equal(t, profileToFeature.User.ID, data.Results[1].ID, "Wrong 2nd profile returned")
				assert.True(t, data.Results[1].IsFeatured, "IsFeatured should be set to true")
			} else {
				assert.Equal(t, currentFeatured.User.ID, data.Results[1].ID, "Wrong 1st profile returned")
				assert.False(t, data.Results[1].IsFeatured, "IsFeatured should not be set")
				assert.Equal(t, profileToFeature.User.ID, data.Results[0].ID, "Wrong 2nd profile returned")
				assert.True(t, data.Results[0].IsFeatured, "IsFeatured should be set to true")
			}
		}
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.BatchUpdate(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	mockTX.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestBatchUpdateFeaturedUserCommitFail(t *testing.T) {
	profileToFeature := testusers.NewProfile()

	currentFeatured := testusers.NewProfile()
	currentFeatured.IsFeatured = ptrs.NewBool(true)

	handlerParams := &users.BatchUpdateParams{
		FeaturedUser: ptrs.NewString(profileToFeature.User.ID),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetID("*users.Profile", profileToFeature.User.ID, func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *profileToFeature
	})
	mockDB.ExpectGetNoParams("*users.Profile", func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *currentFeatured
	})
	mockTX, _ := mockDB.ExpectTransaction()
	mockTX.ExpectUpdate("*users.Profile")
	mockTX.ExpectUpdate("*users.Profile")
	mockTX.ExpectCommitError()
	mockTX.ExpectRollback()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.BatchUpdate(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	mockTX.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}

func TestBatchUpdateFeaturedUserSaveFail(t *testing.T) {
	profileToFeature := testusers.NewProfile()

	currentFeatured := testusers.NewProfile()
	currentFeatured.IsFeatured = ptrs.NewBool(true)

	handlerParams := &users.BatchUpdateParams{
		FeaturedUser: ptrs.NewString(profileToFeature.User.ID),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetID("*users.Profile", profileToFeature.User.ID, func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *profileToFeature
	})
	mockDB.ExpectGetNoParams("*users.Profile", func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *currentFeatured
	})
	mockTX, _ := mockDB.ExpectTransaction()
	mockTX.ExpectUpdate("*users.Profile")
	mockTX.ExpectUpdateError("*users.Profile")
	mockTX.ExpectRollback()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.BatchUpdate(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	mockTX.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}

func TestBatchUpdateFeaturedUserSaveFail2(t *testing.T) {
	profileToFeature := testusers.NewProfile()

	currentFeatured := testusers.NewProfile()
	currentFeatured.IsFeatured = ptrs.NewBool(true)

	handlerParams := &users.BatchUpdateParams{
		FeaturedUser: ptrs.NewString(profileToFeature.User.ID),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetID("*users.Profile", profileToFeature.User.ID, func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *profileToFeature
	})
	mockDB.ExpectGetNoParams("*users.Profile", func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *currentFeatured
	})
	mockTX, _ := mockDB.ExpectTransaction()
	mockTX.ExpectUpdateError("*users.Profile")
	mockTX.ExpectRollback()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.BatchUpdate(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	mockTX.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}

func TestBatchUpdateFeaturedUserAlreadyFeatured(t *testing.T) {
	profileToFeature := testusers.NewProfile()
	profileToFeature.IsFeatured = ptrs.NewBool(true)

	handlerParams := &users.BatchUpdateParams{
		FeaturedUser: ptrs.NewString(profileToFeature.User.ID),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetID("*users.Profile", profileToFeature.User.ID, func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *profileToFeature
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.BatchUpdate(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusConflict, httpErr.HTTPStatus())
	assert.Equal(t, "featured_user", httpErr.Field())
}

func TestBatchUpdateFeaturedUserNotFound(t *testing.T) {
	profileToFeature := testusers.NewProfile()

	handlerParams := &users.BatchUpdateParams{
		FeaturedUser: ptrs.NewString(profileToFeature.User.ID),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetIDNotFound("*users.Profile", profileToFeature.User.ID)

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.BatchUpdate(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, httpErr.HTTPStatus())
}

func TestBatchUpdateFeaturedUserGetFeaturedNotDBCon(t *testing.T) {
	profileToFeature := testusers.NewProfile()

	handlerParams := &users.BatchUpdateParams{
		FeaturedUser: ptrs.NewString(profileToFeature.User.ID),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetID("*users.Profile", profileToFeature.User.ID, func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *profileToFeature
	})
	mockDB.ExpectGetNoParamsError("*users.Profile")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.BatchUpdate(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}

func TestBatchUpdateFeaturedUserTransactionFail(t *testing.T) {
	profileToFeature := testusers.NewProfile()

	currentFeatured := testusers.NewProfile()
	currentFeatured.IsFeatured = ptrs.NewBool(true)

	handlerParams := &users.BatchUpdateParams{
		FeaturedUser: ptrs.NewString(profileToFeature.User.ID),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetID("*users.Profile", profileToFeature.User.ID, func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *profileToFeature
	})
	mockDB.ExpectGetNoParams("*users.Profile", func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *currentFeatured
	})
	mockDB.ExpectTransactionError()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := users.BatchUpdate(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}
