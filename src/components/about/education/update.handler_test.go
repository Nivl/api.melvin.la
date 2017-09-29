package education_test

import (
	"math/rand"
	"net/http"
	"net/url"
	"testing"

	"github.com/melvin-laplanche/ml-api/src/components/about/education/testeducation"

	"github.com/satori/go.uuid"

	"github.com/Nivl/go-params"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-types/ptrs"
	"github.com/melvin-laplanche/ml-api/src/components/about/education"
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
			Description: "Should fail on not nil but empty degree",
			MsgMatch:    params.ErrMsgEmptyParameter,
			FieldName:   "degree",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"degree": []string{"     "},
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
			Description: "Should fail on not nil but empty description",
			MsgMatch:    params.ErrMsgEmptyParameter,
			FieldName:   "description",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"description": []string{"     "},
				},
			},
		},
		{
			Description: "Should fail on invalid start_year type",
			MsgMatch:    params.ErrMsgInvalidInteger,
			FieldName:   "start_year",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"start_year": []string{"not-an-int"},
				},
			},
		},
		{
			Description: "Should fail on invalid start_year",
			MsgMatch:    education.ErrMsgInvalidStartYear,
			FieldName:   "start_year",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"start_year": []string{"5000"},
				},
			},
		},
		{
			Description: "Should fail on invalid end_year",
			MsgMatch:    education.ErrMsgInvalidEndYear,
			FieldName:   "end_year",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"end_year": []string{"1000"},
				},
			},
		},
		{
			Description: "Should fail on invalid end_year type",
			MsgMatch:    params.ErrMsgInvalidInteger,
			FieldName:   "end_year",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"end_year": []string{"not-an-int"},
				},
			},
		},
		{
			Description: "Should fail on end_year before than start year",
			MsgMatch:    education.ErrMsgEndYearBeforeStart,
			FieldName:   "end_year",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"start_year": []string{"2015"},
					"end_year":   []string{"2013"},
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

	g := education.Endpoints[education.EndpointUpdate].Guard
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
			"Should work with a valid in_trash",
			map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"in_trash": []string{"0"},
				},
			},
		},
		{
			"Should work with a empty gpa",
			map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"gpa": []string{""},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := education.Endpoints[education.EndpointUpdate]
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

	g := education.Endpoints[education.EndpointUpdate].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestUpdateHappyPath(t *testing.T) {
	handlerParams := &education.UpdateParams{
		ID:          "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		Degree:      ptrs.NewString("CS"),
		GPA:         ptrs.NewString("4.0"),
		Location:    ptrs.NewString("Location"),
		Description: ptrs.NewString("Description"),
		StartYear:   ptrs.NewInt(rand.Intn(100) + 1950),
		InTrash:     ptrs.NewBool(true),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGet("*education.Education", func(args mock.Arguments) {
		edu := args.Get(0).(*education.Education)
		*edu = *(testeducation.New())
		edu.ID = handlerParams.ID

	})
	mockDB.ExpectUpdate("*education.Education")

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectOk("*education.Payload", func(args mock.Arguments) {
		edu := args.Get(0).(*education.Payload)
		assert.Equal(t, handlerParams.ID, edu.ID, "ID should have not changed")
		assert.Equal(t, *handlerParams.Degree, edu.Degree, "Degree should have been updated")
		assert.Equal(t, *handlerParams.GPA, *edu.GPA, "GPA should have been updated")
		assert.Equal(t, *handlerParams.Location, *edu.Location, "Location should have been updated")
		assert.Equal(t, *handlerParams.Description, *edu.Description, "Description should have been updated")
		assert.Equal(t, *handlerParams.StartYear, edu.StartYear, "StartYear should have been updated")
		assert.NotNil(t, edu.DeletedAt, "DeletedAt should have been set")
		assert.Empty(t, edu.EndYear, "EndYear should have not been set")
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := education.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestUpdateUnsetEndYear(t *testing.T) {
	handlerParams := &education.UpdateParams{
		UnsetEndYear: true,
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGet("*education.Education", func(args mock.Arguments) {
		edu := args.Get(0).(*education.Education)
		edu.ID = "48d0c8b8-d7a3-4855-9d90-29a06ef474b0"
		edu.EndYear = ptrs.NewInt(2017)
	})
	mockDB.ExpectUpdate("*education.Education")

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectOk("*education.Payload", func(args mock.Arguments) {
		edu := args.Get(0).(*education.Payload)
		assert.Empty(t, edu.EndYear, "EndYear should have not been set")
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := education.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestUpdateEndYearBeforeStartYear(t *testing.T) {
	now := rand.Intn(100) + 1950
	future := now + 1

	handlerParams := &education.UpdateParams{
		EndYear: ptrs.NewInt(now),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGet("*education.Education", func(args mock.Arguments) {
		edu := args.Get(0).(*education.Education)
		edu.StartYear = future
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := education.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusBadRequest, apiError.HTTPStatus())
	assert.Equal(t, "end_date", apiError.Field())
}

func TestUpdateNoDBCon(t *testing.T) {
	handlerParams := &education.UpdateParams{
		ID:          "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		Degree:      ptrs.NewString("CS"),
		GPA:         ptrs.NewString("4.0"),
		Location:    ptrs.NewString("Location"),
		Description: ptrs.NewString("Description"),
		StartYear:   ptrs.NewInt(rand.Intn(100) + 1950),
		InTrash:     ptrs.NewBool(true),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGet("*education.Education", func(args mock.Arguments) {
		exp := args.Get(0).(*education.Education)
		*exp = *(testeducation.New())
	})
	mockDB.ExpectUpdateError("*education.Education")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := education.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, apiError.HTTPStatus())
}

func TestUpdateUnexisting(t *testing.T) {
	handlerParams := &education.UpdateParams{
		ID:          "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		Degree:      ptrs.NewString("CS"),
		GPA:         ptrs.NewString("4.0"),
		Location:    ptrs.NewString("Location"),
		Description: ptrs.NewString("Description"),
		StartYear:   ptrs.NewInt(rand.Intn(100) + 1950),
		InTrash:     ptrs.NewBool(true),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetNotFound("*education.Education")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := education.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, apiError.HTTPStatus())
}

func TestUpdateUnexistingOrg(t *testing.T) {
	handlerParams := &education.UpdateParams{
		ID:             uuid.NewV4().String(),
		OrganizationID: ptrs.NewString(uuid.NewV4().String()),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGet("*education.Education", func(args mock.Arguments) {
		exp := args.Get(0).(*education.Education)
		*exp = *(testeducation.New())
	})
	mockDB.ExpectGetNotFound("*organizations.Organization")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := education.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, apiError.HTTPStatus())
	assert.Equal(t, "organization_id", apiError.Field())
}

func TestUpdateGetOrgNoDBCon(t *testing.T) {
	handlerParams := &education.UpdateParams{
		ID:             uuid.NewV4().String(),
		OrganizationID: ptrs.NewString(uuid.NewV4().String()),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGet("*education.Education", func(args mock.Arguments) {
		exp := args.Get(0).(*education.Education)
		*exp = *(testeducation.New())
	})
	mockDB.ExpectGetError("*organizations.Organization")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := education.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, apiError.HTTPStatus())
}
