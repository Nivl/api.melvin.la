package organizations_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-rest-tools/primitives/ptrs"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateInvalidParams(t *testing.T) {
	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing ID",
			MsgMatch:    "parameter missing",
			FieldName:   "id",
			Sources: map[string]url.Values{
				"url": url.Values{},
			},
		},
		{
			Description: "Should fail on invalid ID",
			MsgMatch:    "not a valid uuid",
			FieldName:   "id",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"xxx"},
				},
			},
		},
		{
			Description: "Should fail on not nil but empty name",
			MsgMatch:    "parameter can be omitted but not empty",
			FieldName:   "name",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"name": []string{"     "},
				},
			},
		},
		{
			Description: "Should fail on not nil but invalid website",
			MsgMatch:    "not a valid url",
			FieldName:   "website",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"name":    []string{"valid name"},
					"website": []string{"not-a-url"},
				},
			},
		},
		{
			Description: "Should fail on not nil but invalid in_trash",
			MsgMatch:    "invalid boolean",
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
	}

	g := organizations.Endpoints[organizations.EndpointUpdate].Guard
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
			"Should work with only a valid name",
			map[string]url.Values{
				"url": url.Values{
					"id": []string{"aa44ca86-553e-4e16-8c30-2e50e63f7eaa"},
				},
				"form": url.Values{
					"name": []string{"valid name"},
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

			endpts := organizations.Endpoints[organizations.EndpointUpdate]
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

	g := organizations.Endpoints[organizations.EndpointUpdate].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestUpdateHappyPath(t *testing.T) {
	handlerParams := &organizations.UpdateParams{
		ID:        "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		Name:      ptrs.NewString("Google"),
		ShortName: ptrs.NewString("googl"),
		Website:   ptrs.NewString("https://google.com"),
		InTrash:   ptrs.NewBool(true),
	}

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGet("*organizations.Organization", func(args mock.Arguments) {
		org := args.Get(0).(*organizations.Organization)
		org.ID = "48d0c8b8-d7a3-4855-9d90-29a06ef474b0"
		org.Name = "Not Google"
	})
	mockDB.ExpectUpdate("*organizations.Organization")

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectOk("*organizations.Organization", func(args mock.Arguments) {
		org := args.Get(0).(*organizations.Organization)
		assert.Equal(t, handlerParams.ID, org.ID, "ID should have not changed")
		assert.Equal(t, *handlerParams.Name, org.Name, "name should have been updated")
		assert.Equal(t, *handlerParams.Website, *org.Website, "Website should have been updated")
		assert.Equal(t, *handlerParams.ShortName, *org.ShortName, "ShortName should have been updated")
		assert.NotNil(t, *org.DeletedAt, "DeletedAt should have been set")
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := organizations.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Nil(t, err, "the handler should not have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

// func TestUpdateConflictName(t *testing.T) {
// 	p := &testrouter.ConflictTestParams{
// 		StructConflicting: "*organizations.Organization",
// 		FieldConflicting:  "name",
// 		Handler:           organizations.Update,
// 		HandlerParams: &organizations.UpdateParams{
// 			Name: ptrs.NewString("Google"),
// 		},
// 	}
// 	testrouter.ConflictUpdateTest(t, p)
// }

// func TestUpdateConflictShortName(t *testing.T) {
// 	p := &testrouter.ConflictTestParams{
// 		StructConflicting: "*organizations.Organization",
// 		FieldConflicting:  "short_name",
// 		Handler:           organizations.Update,
// 		HandlerParams: &organizations.UpdateParams{
// 			Name:      "Google",
// 			ShortName: ptrs.NewString("googl"),
// 		},
// 	}
// 	testrouter.ConflictUpdateTest(t, p)
// }
