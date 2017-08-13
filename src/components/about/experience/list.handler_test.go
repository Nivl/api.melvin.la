package experience_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/params"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience"
	"github.com/stretchr/testify/assert"
)

func TestListInvalidParams(t *testing.T) {
	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on invalid delete value",
			MsgMatch:    params.ErrMsgInvalidBoolean,
			FieldName:   "deleted",
			Sources: map[string]url.Values{
				"query": url.Values{
					"deleted": []string{"not a bool"},
				},
			},
		},
		{
			Description: "Should fail on invalid orphans value",
			MsgMatch:    params.ErrMsgInvalidBoolean,
			FieldName:   "orphans",
			Sources: map[string]url.Values{
				"query": url.Values{
					"orphans": []string{"not a bool"},
				},
			},
		},
		{
			Description: "Should fail with page = 0",
			MsgMatch:    paginator.ErrMsgNumberBelow1,
			FieldName:   "page",
			Sources: map[string]url.Values{
				"query": url.Values{
					"page": []string{"0"},
				},
			},
		},
		{
			Description: "Should fail with per_page = 0",
			MsgMatch:    paginator.ErrMsgNumberBelow1,
			FieldName:   "per_page",
			Sources: map[string]url.Values{
				"query": url.Values{
					"per_page": []string{"0"},
				},
			},
		},
		{
			Description: "Should fail with per_page > 100",
			MsgMatch:    "cannot be > 100",
			FieldName:   "per_page",
			Sources: map[string]url.Values{
				"query": url.Values{
					"per_page": []string{"101"},
				},
			},
		},
		{
			Description: "Should fail with invalid operator",
			MsgMatch:    params.ErrMsgEnum,
			FieldName:   "op",
			Sources: map[string]url.Values{
				"query": url.Values{
					"op": []string{"nand"},
				},
			},
		},
	}

	g := experience.Endpoints[experience.EndpointList].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestListValidParams(t *testing.T) {
	testCases := []struct {
		description string
		sources     map[string]url.Values
	}{
		{
			"Should work with nothing",
			map[string]url.Values{
				"query": url.Values{},
			},
		},
		{
			"Should work with deleted=true",
			map[string]url.Values{
				"query": url.Values{
					"deleted": []string{"true"},
				},
			},
		},
		{
			"Should work with deleted=false",
			map[string]url.Values{
				"query": url.Values{
					"deleted": []string{"false"},
				},
			},
		},
		{
			"Should work with orphans=true",
			map[string]url.Values{
				"query": url.Values{
					"orphans": []string{"true"},
				},
			},
		},
		{
			"Should work with orphans=false",
			map[string]url.Values{
				"query": url.Values{
					"orphans": []string{"false"},
				},
			},
		},
		{
			"Should work with both orphans and deleted",
			map[string]url.Values{
				"query": url.Values{
					"orphans": []string{"true"},
					"deleted": []string{"true"},
				},
			},
		},
		{
			"Should work with orphans, deleted and an op",
			map[string]url.Values{
				"query": url.Values{
					"orphans": []string{"true"},
					"deleted": []string{"true"},
					"op":      []string{"or"},
				},
			},
		},
		{
			"Should work with a page",
			map[string]url.Values{
				"query": url.Values{
					"page": []string{"1"},
				},
			},
		},
		{
			"Should work with a per_page",
			map[string]url.Values{
				"query": url.Values{
					"per_page": []string{"10"},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := experience.Endpoints[experience.EndpointList]
			_, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)
		})
	}
}

func TestListAccess(t *testing.T) {
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

	g := experience.Endpoints[experience.EndpointList].Guard
	testguard.AccessTest(t, g, testCases)
}
