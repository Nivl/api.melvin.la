package experience_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-params"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-types/date"
	"github.com/dchest/uniuri"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
)

func TestAddInvalidParams(t *testing.T) {
	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing org ID",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "organization_id",
			Sources: map[string]url.Values{
				"form": url.Values{
					"job_title":   []string{"job Title"},
					"location":    []string{"Los Angeles, CA"},
					"description": []string{"description of the work done"},
					"start_date":  []string{"2016-05"},
				},
			},
		},
		{
			Description: "Should fail on invalid org ID",
			MsgMatch:    params.ErrMsgInvalidUUID,
			FieldName:   "organization_id",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"not-a-uuid"},
					"job_title":       []string{"job Title"},
					"location":        []string{"Los Angeles, CA"},
					"description":     []string{"description of the work done"},
					"start_date":      []string{"2016-05"},
				},
			},
		},
		{
			Description: "Should fail on missing Job Title",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "job_title",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"location":        []string{"Los Angeles, CA"},
					"description":     []string{"description of the work done"},
					"start_date":      []string{"2016-05"},
				},
			},
		},
		{
			Description: "Should fail on invalid Job Title",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "job_title",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"job_title":       []string{"     "},
					"location":        []string{"Los Angeles, CA"},
					"description":     []string{"description of the work done"},
					"start_date":      []string{"2016-05"},
				},
			},
		},
		{
			Description: "Should fail on long Job Title",
			MsgMatch:    params.ErrMsgMaxLen,
			FieldName:   "job_title",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"job_title":       []string{uniuri.NewLen(256)},
					"location":        []string{"Los Angeles, CA"},
					"description":     []string{"description of the work done"},
					"start_date":      []string{"2016-05"},
				},
			},
		},
		{
			Description: "Should fail on missing location",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "location",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"job_title":       []string{"job Title"},
					"description":     []string{"description of the work done"},
					"start_date":      []string{"2016-05"},
				},
			},
		},
		{
			Description: "Should fail on invalid location",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "location",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"job_title":       []string{"job Title"},
					"location":        []string{"     "},
					"description":     []string{"description of the work done"},
					"start_date":      []string{"2016-05"},
				},
			},
		},
		{
			Description: "Should fail on long location",
			MsgMatch:    params.ErrMsgMaxLen,
			FieldName:   "location",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"job_title":       []string{"job Title"},
					"location":        []string{uniuri.NewLen(256)},
					"description":     []string{"description of the work done"},
					"start_date":      []string{"2016-05"},
				},
			},
		},
		{
			Description: "Should fail on missing description",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "description",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"job_title":       []string{"job Title"},
					"location":        []string{"Los Angeles, CA"},
					"start_date":      []string{"2016-05"},
				},
			},
		},
		{
			Description: "Should fail on missing description",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "description",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"job_title":       []string{"job Title"},
					"location":        []string{"Los Angeles, CA"},
					"description":     []string{"      "},
					"start_date":      []string{"2016-05"},
				},
			},
		},
		{
			Description: "Should fail on long description",
			MsgMatch:    params.ErrMsgMaxLen,
			FieldName:   "description",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"job_title":       []string{"job Title"},
					"location":        []string{"Los Angeles, CA"},
					"description":     []string{uniuri.NewLen(10001)},
					"start_date":      []string{"2016-05"},
				},
			},
		},
		{
			Description: "Should fail on missing start date",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "start_date",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"job_title":       []string{"job Title"},
					"location":        []string{"Los Angeles, CA"},
					"description":     []string{"description of the work done"},
				},
			},
		},
		{
			Description: "Should fail on invalid start date",
			MsgMatch:    date.ErrMsgInvalidFormat,
			FieldName:   "start_date",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"job_title":       []string{"job Title"},
					"location":        []string{"Los Angeles, CA"},
					"description":     []string{"description of the work done"},
					"start_date":      []string{"05-2016"},
				},
			},
		},
		{
			Description: "Should fail with the end date being before the start date",
			MsgMatch:    experience.ErrMsgInvalidEndDate,
			FieldName:   "end_date",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"job_title":       []string{"job Title"},
					"location":        []string{"Los Angeles, CA"},
					"description":     []string{"description of the work done"},
					"start_date":      []string{"2016-05"},
					"end_date":        []string{"2015-04"},
				},
			},
		},
	}

	g := experience.Endpoints[experience.EndpointAdd].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestAddValidParams(t *testing.T) {
	testCases := []struct {
		description string
		sources     map[string]url.Values
	}{
		{
			"Should work with all the fields",
			map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"job_title":       []string{"job Title"},
					"location":        []string{"Los Angeles, CA"},
					"description":     []string{"description of the work done"},
					"start_date":      []string{"2016-05"},
					"end_date":        []string{"2017-05"},
				},
			},
		},
		{
			"Should work without end_date",
			map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"job_title":       []string{"job Title"},
					"location":        []string{"Los Angeles, CA"},
					"description":     []string{"description of the work done"},
					"start_date":      []string{"2016-05"},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := experience.Endpoints[experience.EndpointAdd]
			_, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)
		})
	}
}

func TestAddAccess(t *testing.T) {
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

	g := experience.Endpoints[experience.EndpointAdd].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestAddHappyPath(t *testing.T) {
	org := &organizations.Organization{
		ID:   "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		Name: "Org name",
	}

	handlerParams := &experience.AddParams{
		OrganizationID: org.ID,
		JobTitle:       "Title",
		Location:       "Los Angeles area, CA",
		Description:    "description of the work done",
		StartDate:      date.Today(),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectInsert("*experience.Experience")
	mockDB.ExpectGet("*organizations.Organization", func(args mock.Arguments) {
		o := args.Get(0).(*organizations.Organization)
		*o = *org
	})

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectCreated("*experience.Payload", func(args mock.Arguments) {
		exp := args.Get(0).(*experience.Payload)
		assert.Equal(t, handlerParams.JobTitle, exp.JobTitle)
		assert.Equal(t, handlerParams.Location, exp.Location)
		assert.Equal(t, handlerParams.Description, exp.Description)
		assert.Equal(t, org.ID, exp.Organization.ID)
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := experience.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestAddOrgNotFound(t *testing.T) {
	handlerParams := &experience.AddParams{
		OrganizationID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		JobTitle:       "Title",
		Location:       "Los Angeles area, CA",
		Description:    "description of the work done",
		StartDate:      date.Today(),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetNotFound("*organizations.Organization")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := experience.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, httpErr.HTTPStatus())
}

func TestAddNoDBCon(t *testing.T) {
	org := &organizations.Organization{
		ID:   "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		Name: "Org name",
	}

	handlerParams := &experience.AddParams{
		OrganizationID: org.ID,
		JobTitle:       "Title",
		Location:       "Los Angeles area, CA",
		Description:    "description of the work done",
		StartDate:      date.Today(),
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectInsertError("*experience.Experience")
	mockDB.ExpectGet("*organizations.Organization", func(args mock.Arguments) {
		o := args.Get(0).(*organizations.Organization)
		*o = *org
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := experience.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}
