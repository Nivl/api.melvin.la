package organizations_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-rest-tools/primitives/ptrs"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/router/testrouter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInvalidParams(t *testing.T) {
	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing name",
			MsgMatch:    "parameter missing",
			FieldName:   "name",
			Sources: map[string]url.Values{
				"form": url.Values{
					"short_name": []string{"short name"},
					"website":    []string{"http://domain.tld"},
				},
			},
		},
		{
			Description: "Should fail on invalid name",
			MsgMatch:    "parameter missing",
			FieldName:   "name",
			Sources: map[string]url.Values{
				"form": url.Values{
					"name": []string{"     "},
				},
			},
		},
		{
			Description: "Should fail on invalid website",
			MsgMatch:    "not a valid url",
			FieldName:   "website",
			Sources: map[string]url.Values{
				"form": url.Values{
					"name":    []string{"valid name"},
					"website": []string{"not-a-url"},
				},
			},
		},
	}

	g := organizations.Endpoints[organizations.EndpointAdd].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestAddValidParams(t *testing.T) {
	testCases := []struct {
		description string
		sources     map[string]url.Values
	}{
		{
			"Should work with only a valid name",
			map[string]url.Values{
				"form": url.Values{
					"name": []string{"name"},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := organizations.Endpoints[organizations.EndpointAdd]
			data, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)

			if data != nil {
				p := data.(*organizations.AddParams)
				assert.Equal(t, tc.sources["form"].Get("name"), p.Name)
			}
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

	g := organizations.Endpoints[organizations.EndpointAdd].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestAddHappyPath(t *testing.T) {
	handlerParams := &organizations.AddParams{
		Name:    "Google",
		Website: ptrs.NewString("https://google.com"),
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectInsert("*organizations.Organization")

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectCreated("*organizations.Organization", func(args mock.Arguments) {
		org := args.Get(0).(*organizations.Organization)
		assert.Equal(t, handlerParams.Name, org.Name)
		assert.Equal(t, *handlerParams.Website, *org.Website)
		assert.NotEmpty(t, org.ID)
		assert.Empty(t, org.ShortName)
		assert.Empty(t, org.Logo)
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := organizations.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Nil(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestAddConflictName(t *testing.T) {
	p := &testrouter.ConflictTestParams{
		StructConflicting: "*organizations.Organization",
		FieldConflicting:  "name",
		Handler:           organizations.Add,
		HandlerParams: &organizations.AddParams{
			Name: "Google",
		},
	}
	testrouter.ConflictTest(t, p)
}

func TestAddConflictShortName(t *testing.T) {
	p := &testrouter.ConflictTestParams{
		StructConflicting: "*organizations.Organization",
		FieldConflicting:  "short_name",
		Handler:           organizations.Add,
		HandlerParams: &organizations.AddParams{
			Name:      "Google",
			ShortName: ptrs.NewString("googl"),
		},
	}
	testrouter.ConflictTest(t, p)
}
