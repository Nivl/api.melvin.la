package organizations

import (
	"fmt"

	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
)

var listEndpoint = &router.Endpoint{
	Verb:    "GET",
	Path:    "/about/organizations",
	Handler: List,
	Guard: &guard.Guard{
		Auth:        guard.AdminAccess,
		ParamStruct: &ListParams{},
	},
}

// ListParams represents the params accepted by the Add endpoint
type ListParams struct {
	paginator.HandlerParams
	Deleted bool `from:"query" json:"deleted" default:"false"`
}

// List is an endpoint used to list all Organization
func List(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*ListParams)
	paginator := params.Paginator()

	whereClause := ""
	if !params.Deleted {
		whereClause = "WHERE deleted_at IS NULL"
	}

	orgs := Organizations{}
	stmt := fmt.Sprintf(`SELECT * from about_organizations %s ORDER BY name OFFSET $1 LIMIT $2`, whereClause)

	err := deps.DB.Select(&orgs, stmt, paginator.Offset(), paginator.Limit())
	if err != nil {
		return err
	}
	return req.Response().Ok(orgs.Export())
}
