package organizations

import (
	"net/http"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
)

var addEndpoint = &router.Endpoint{
	Verb:    http.MethodPost,
	Path:    "/about/organizations",
	Handler: Add,
	Guard: &guard.Guard{
		Auth:        guard.AdminAccess,
		ParamStruct: &AddParams{},
	},
}

// AddParams represents the params accepted by the Add endpoint
type AddParams struct {
	Name      string  `from:"form" json:"name" params:"required,trim"`
	ShortName *string `from:"form" json:"short_name" params:"trim"`
	Website   *string `from:"form" json:"website" params:"url,trim"`
}

// Add is an endpoint used to add a new Organization
func Add(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*AddParams)

	org := &Organization{
		Name: params.Name,
	}

	if params.Website != nil && *params.Website != "" {
		org.Website = params.Website
	}
	if params.ShortName != nil && *params.ShortName != "" {
		org.ShortName = params.ShortName
	}

	if err := org.Create(deps.DB); err != nil {
		return err
	}
	return req.Response().Created(org.ExportPrivate())
}
