package experience_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/satori/go.uuid"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/router/params"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience/testexperience"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetInvalidParams(t *testing.T) {
	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing ID",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "id",
			Sources: map[string]url.Values{
				"url": url.Values{},
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
			},
		},
	}

	g := experience.Endpoints[experience.EndpointGet].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestGetValidParams(t *testing.T) {
	testCases := []struct {
		description string
		sources     map[string]url.Values
	}{
		{
			"Should work with a valid ID",
			map[string]url.Values{
				"url": url.Values{
					"id": []string{"48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := users.Endpoints[users.EndpointGet]
			data, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)

			if data != nil {
				p := data.(*users.GetParams)
				assert.Equal(t, tc.sources["url"].Get("id"), p.ID)
			}
		})
	}
}

func TestGetHappyPath(t *testing.T) {
	testCases := []struct {
		description string
		requester   *auth.User
	}{
		{"anonynous user", nil},
		{"logged-in user", testauth.NewUser()},
		{"admin user", testauth.NewAdmin()},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			exp := testexperience.New()
			handlerParams := &experience.GetParams{
				ID: exp.ID,
			}

			// Mock the database & add expectations
			mockDB := new(mockdb.DB)
			mockDB.ExpectGet("*experience.Experience", func(args mock.Arguments) {
				data := args.Get(0).(*experience.Experience)
				*data = *exp
			})

			// Mock the response & add expectations
			res := new(mockrouter.HTTPResponse)
			res.ExpectOk("*experience.Payload", func(args mock.Arguments) {
				pld := args.Get(0).(*experience.Payload)
				assert.Equal(t, exp.ID, pld.ID, "ID should not have changed")
				assert.Equal(t, exp.JobTitle, pld.JobTitle, "Name should not have changed")
				assert.Equal(t, exp.Description, pld.Description, "the description should not have changed")
				assert.Equal(t, exp.Location, pld.Location, "the location should not have changed")
				assert.Equal(t, exp.OrganizationID, pld.Organization.ID, "OrganizationID should not have changed")
				assert.Equal(t, exp.Organization.ID, pld.Organization.ID, "Organization.ID should not have changed")
				assert.Equal(t, exp.Organization.Name, pld.Organization.Name, "Organization Name should not have changed")

				if tc.requester.IsAdm() {
					assert.NotNil(t, pld.CreatedAt, "CreatedAt should have been returned")
					assert.NotNil(t, pld.UpdatedAt, "UpdatedAt should have been returned")
					assert.Nil(t, pld.DeletedAt, "UpdatedAt should not have been set")
					assert.NotNil(t, pld.Organization.CreatedAt, "Organization's CreatedAt should have been set")
					assert.NotNil(t, pld.Organization.UpdatedAt, "Organization's CreatedAt should have been set")
				} else {
					assert.Nil(t, pld.CreatedAt, "CreatedAt should not have been returned")
					assert.Nil(t, pld.UpdatedAt, "UpdatedAt should not have been returned")
					assert.Nil(t, pld.DeletedAt, "UpdatedAt should not have been returned nor set")
					assert.Nil(t, pld.Organization.CreatedAt, "Organization's CreatedAt should not have been set")
					assert.Nil(t, pld.Organization.UpdatedAt, "Organization's CreatedAt should not have been set")
				}
			})

			// Mock the request & add expectations
			req := new(mockrouter.HTTPRequest)
			req.On("Response").Return(res)
			req.On("Params").Return(handlerParams)
			req.On("User").Return(tc.requester)

			// call the handler
			err := experience.Get(req, &router.Dependencies{DB: mockDB})

			// Assert everything
			assert.NoError(t, err, "the handler should not have fail")
			mockDB.AssertExpectations(t)
			req.AssertExpectations(t)
			res.AssertExpectations(t)
		})
	}
}

func TestGetUnexisting(t *testing.T) {
	handlerParams := &experience.GetParams{
		ID: uuid.NewV4().String(),
	}
	requester := testauth.NewUser()

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGetNotFound("*experience.Experience")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(requester)

	// call the handler
	err := experience.Get(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, httpErr.HTTPStatus())
}
