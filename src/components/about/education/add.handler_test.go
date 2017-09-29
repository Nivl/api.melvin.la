package education_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-params"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/dchest/uniuri"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/melvin-laplanche/ml-api/src/components/about/education"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations/testorganizations"
)

func TestAddInvalidParams(t *testing.T) {
	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing org ID",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "organization_id",
			Sources: map[string]url.Values{
				"form": url.Values{
					"degree":      []string{"CS"},
					"location":    []string{"CSULB"},
					"description": []string{"description"},
					"start_year":  []string{"2013"},
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
					"degree":          []string{"CS"},
					"location":        []string{"CSULB"},
					"description":     []string{"description"},
					"start_year":      []string{"2013"},
				},
			},
		},
		{
			Description: "Should fail on missing Degree",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "degree",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"location":        []string{"CSULB"},
					"description":     []string{"description"},
					"start_year":      []string{"2013"},
				},
			},
		},
		{
			Description: "Should fail on invalid Degree",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "degree",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"degree":          []string{"    "},
					"location":        []string{"CSULB"},
					"description":     []string{"description"},
					"start_year":      []string{"2013"},
				},
			},
		},
		{
			Description: "Should fail on long degree",
			MsgMatch:    params.ErrMsgMaxLen,
			FieldName:   "degree",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"degree":          []string{uniuri.NewLen(256)},
					"location":        []string{"CSULB"},
					"description":     []string{"description"},
					"start_year":      []string{"2013"},
				},
			},
		},
		{
			Description: "Should fail on long gpa",
			MsgMatch:    params.ErrMsgMaxLen,
			FieldName:   "gpa",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"gpa":             []string{uniuri.NewLen(6)},
					"degree":          []string{"CS"},
					"location":        []string{"CSULB"},
					"description":     []string{"description"},
					"start_year":      []string{"2013"},
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
					"degree":          []string{"CS"},
					"description":     []string{"description"},
					"start_year":      []string{"2013"},
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
					"degree":          []string{"CS"},
					"location":        []string{"      "},
					"description":     []string{"description"},
					"start_year":      []string{"2013"},
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
					"degree":          []string{"CS"},
					"location":        []string{uniuri.NewLen(256)},
					"description":     []string{"description"},
					"start_year":      []string{"2013"},
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
					"degree":          []string{"CS"},
					"location":        []string{"CSULB"},
					"start_year":      []string{"2013"},
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
					"degree":          []string{"CS"},
					"location":        []string{"CSULB"},
					"description":     []string{"     "},
					"start_year":      []string{"2013"},
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
					"degree":          []string{"CS"},
					"location":        []string{"CSULB"},
					"description":     []string{uniuri.NewLen(10001)},
					"start_year":      []string{"2013"},
				},
			},
		},
		{
			Description: "Should fail on missing start year",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "start_year",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"degree":          []string{"CS"},
					"location":        []string{"CSULB"},
					"description":     []string{"description"},
				},
			},
		},
		{
			Description: "Should fail on invalid start year",
			MsgMatch:    education.ErrMsgInvalidStartYear,
			FieldName:   "start_year",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"degree":          []string{"CS"},
					"location":        []string{"CSULB"},
					"description":     []string{"description"},
					"start_year":      []string{"1"},
				},
			},
		},
		{
			Description: "Should fail with the end year being before the start year",
			MsgMatch:    education.ErrMsgEndYearBeforeStart,
			FieldName:   "end_year",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"degree":          []string{"CS"},
					"location":        []string{"CSULB"},
					"description":     []string{"description"},
					"start_year":      []string{"2013"},
					"end_year":        []string{"2011"},
				},
			},
		},
		{
			Description: "Should fail on invalid end year",
			MsgMatch:    education.ErrMsgInvalidEndYear,
			FieldName:   "end_year",
			Sources: map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"degree":          []string{"CS"},
					"location":        []string{"CSULB"},
					"description":     []string{"description"},
					"start_year":      []string{"2013"},
					"end_year":        []string{"3098"},
				},
			},
		},
	}

	g := education.Endpoints[education.EndpointAdd].Guard
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
					"degree":          []string{"CS"},
					"gpa":             []string{"4.0"},
					"location":        []string{"CSULB"},
					"description":     []string{"description"},
					"start_year":      []string{"2013"},
					"end_year":        []string{"2013"},
				},
			},
		},
		{
			"Should work without end_date and gpa",
			map[string]url.Values{
				"form": url.Values{
					"organization_id": []string{"7cfb9bd6-7e5d-4793-b0bd-d0a26a758390"},
					"degree":          []string{"CS"},
					"location":        []string{"CSULB"},
					"description":     []string{"description"},
					"start_year":      []string{"2013"},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := education.Endpoints[education.EndpointAdd]
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

	g := education.Endpoints[education.EndpointAdd].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestAddHappyPath(t *testing.T) {
	org := testorganizations.New()

	handlerParams := &education.AddParams{
		OrganizationID: org.ID,
		Degree:         "CS",
		Location:       "Long Beach, CA",
		Description:    "description",
		StartYear:      2013,
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectInsert("*education.Education")
	mockDB.ExpectGet("*organizations.Organization", func(args mock.Arguments) {
		o := args.Get(0).(*organizations.Organization)
		*o = *org
	})

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectCreated("*education.Payload", func(args mock.Arguments) {
		edu := args.Get(0).(*education.Payload)
		assert.Equal(t, handlerParams.Degree, edu.Degree)
		assert.Equal(t, handlerParams.Location, *edu.Location)
		assert.Equal(t, handlerParams.Description, *edu.Description)
		assert.Equal(t, org.ID, edu.Organization.ID)
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := education.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestAddOrgNotFound(t *testing.T) {
	handlerParams := &education.AddParams{
		OrganizationID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		Degree:         "Title",
		Location:       "Long Beach, CA",
		Description:    "description",
		StartYear:      2013,
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetNotFound("*organizations.Organization")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := education.Add(req, &router.Dependencies{DB: mockDB})

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

	handlerParams := &education.AddParams{
		OrganizationID: org.ID,
		Degree:         "Title",
		Location:       "Long Beach, CA",
		Description:    "description",
		StartYear:      2013,
	}

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectInsertError("*education.Education")
	mockDB.ExpectGet("*organizations.Organization", func(args mock.Arguments) {
		o := args.Get(0).(*organizations.Organization)
		*o = *org
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := education.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}
