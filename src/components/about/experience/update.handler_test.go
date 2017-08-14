package experience_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/melvin-laplanche/ml-api/src/components/about/experience/testexperience"

	"github.com/satori/go.uuid"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/router/params"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-rest-tools/types/ptrs"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateInvalidParams(t *testing.T) {
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
					"id": []string{"xxx"},
				},
			},
		},
		{
			Description: "Should fail on not nil but empty job title",
			MsgMatch:    params.ErrMsgEmptyParameter,
			FieldName:   "job_title",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"job_title": []string{"     "},
				},
			},
		},
		{
			Description: "Should fail on not nil but empty location",
			MsgMatch:    params.ErrMsgEmptyParameter,
			FieldName:   "location",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"location": []string{"     "},
				},
			},
		},
		{
			Description: "Should fail on invalid start_date",
			MsgMatch:    params.ErrMsgInvalidDate,
			FieldName:   "start_date",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"start_date": []string{"xxxx-yy"},
				},
			},
		},
		{
			Description: "Should fail on invalid end_date",
			MsgMatch:    params.ErrMsgInvalidDate,
			FieldName:   "end_date",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"end_date": []string{"whatever"},
				},
			},
		},
		{
			Description: "Should fail on not nil but invalid in_trash",
			MsgMatch:    params.ErrMsgInvalidBoolean,
			FieldName:   "in_trash",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"in_trash": []string{"not-a-boolean"},
				},
			},
		},
		{
			Description: "Should fail on not nil but invalid organization_id",
			MsgMatch:    params.ErrMsgInvalidUUID,
			FieldName:   "organization_id",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"organization_id": []string{"not-a-uuid"},
				},
			},
		},
	}

	g := experience.Endpoints[experience.EndpointUpdate].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestUpdateValidParams(t *testing.T) {
	testCases := []struct {
		description string
		sources     map[string]url.Values
	}{
		{
			"Should work with only a valid uuid",
			map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{},
			},
		},
		{
			"Should work with only a valid job_title",
			map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"job_title": []string{"valid name"},
				},
			},
		},
		{
			"Should work with only a valid in_trash",
			map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"in_trash": []string{"0"},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := experience.Endpoints[experience.EndpointUpdate]
			_, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)
		})
	}
}

func TestUpdateAccess(t *testing.T) {
	testCases := []testguard.AccessTestCase{
		{
			Description: "Should fail for anonymous users",
			User:        nil,
			ErrCode:     http.StatusUnauthorized,
		},
		{
			Description: "Should fail for logged users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
			ErrCode:     http.StatusForbidden,
		},
		{
			Description: "Should work for admin users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0", IsAdmin: true},
			ErrCode:     0,
		},
	}

	g := experience.Endpoints[experience.EndpointUpdate].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestUpdateHappyPath(t *testing.T) {
	handlerParams := &experience.UpdateParams{
		ID:          "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		JobTitle:    ptrs.NewString("JobTitle"),
		Location:    ptrs.NewString("Location"),
		Description: ptrs.NewString("Description"),
		StartDate:   db.Today(),
		InTrash:     ptrs.NewBool(true),
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGet("*experience.Experience", func(args mock.Arguments) {
		exp := args.Get(0).(*experience.Experience)
		exp.ID = "48d0c8b8-d7a3-4855-9d90-29a06ef474b0"
		exp.JobTitle = "old JobTitle"
		exp.Location = "old Location"
		exp.Description = "old Description"
		exp.StartDate, _ = db.NewDate("2016-01")
	})
	mockDB.ExpectUpdate("*experience.Experience")

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectOk("*experience.Payload", func(args mock.Arguments) {
		exp := args.Get(0).(*experience.Payload)
		assert.Equal(t, handlerParams.ID, exp.ID, "ID should have not changed")
		assert.Equal(t, *handlerParams.JobTitle, exp.JobTitle, "JobTitle should have been updated")
		assert.Equal(t, *handlerParams.Location, exp.Location, "Location should have been updated")
		assert.Equal(t, *handlerParams.Description, exp.Description, "Description should have been updated")
		assert.Equal(t, handlerParams.StartDate.String(), exp.StartDate.String(), "StartDate should have been updated")
		assert.NotNil(t, exp.DeletedAt, "DeletedAt should have been set")
		assert.Nil(t, exp.EndDate, "EndDate should have not been set")
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := experience.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestUpdateUnsetEndDate(t *testing.T) {
	handlerParams := &experience.UpdateParams{
		UnsetEndDate: true,
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGet("*experience.Experience", func(args mock.Arguments) {
		exp := args.Get(0).(*experience.Experience)
		exp.ID = "48d0c8b8-d7a3-4855-9d90-29a06ef474b0"
		exp.EndDate = db.Today()
	})
	mockDB.ExpectUpdate("*experience.Experience")

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectOk("*experience.Payload", func(args mock.Arguments) {
		exp := args.Get(0).(*experience.Payload)
		assert.Nil(t, exp.EndDate, "EndDate should have not been set")
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := experience.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestUpdateEndDateBeforeStartDate(t *testing.T) {
	now, _ := db.NewDate("2017-08")
	future, _ := db.NewDate("2018-08")

	handlerParams := &experience.UpdateParams{
		EndDate: now,
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGet("*experience.Experience", func(args mock.Arguments) {
		exp := args.Get(0).(*experience.Experience)
		exp.StartDate = future
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := experience.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusBadRequest, apiError.HTTPStatus())
	assert.Equal(t, "end_date", apiError.Field())
}

func TestUpdateNoDBCon(t *testing.T) {
	handlerParams := &experience.UpdateParams{
		ID:          "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		JobTitle:    ptrs.NewString("JobTitle"),
		Location:    ptrs.NewString("Location"),
		Description: ptrs.NewString("Description"),
		StartDate:   db.Today(),
		InTrash:     ptrs.NewBool(true),
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGet("*experience.Experience", func(args mock.Arguments) {
		exp := args.Get(0).(*experience.Experience)
		*exp = *(testexperience.New())
	})
	mockDB.ExpectUpdateError("*experience.Experience")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := experience.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, apiError.HTTPStatus())
}

func TestUpdateUnexisting(t *testing.T) {
	handlerParams := &experience.UpdateParams{
		ID:          "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		JobTitle:    ptrs.NewString("JobTitle"),
		Location:    ptrs.NewString("Location"),
		Description: ptrs.NewString("Description"),
		StartDate:   db.Today(),
		InTrash:     ptrs.NewBool(true),
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGetNotFound("*experience.Experience")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := experience.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, apiError.HTTPStatus())
}

func TestUpdateUnexistingOrg(t *testing.T) {
	handlerParams := &experience.UpdateParams{
		ID:             uuid.NewV4().String(),
		OrganizationID: ptrs.NewString(uuid.NewV4().String()),
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGet("*experience.Experience", func(args mock.Arguments) {
		exp := args.Get(0).(*experience.Experience)
		*exp = *(testexperience.New())
	})
	mockDB.ExpectGetNotFound("*organizations.Organization")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := experience.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, apiError.HTTPStatus())
	assert.Equal(t, "organization_id", apiError.Field())
}
