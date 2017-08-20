package education

import (
	"strings"

	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
)

var listEndpoint = &router.Endpoint{
	Verb:    "GET",
	Path:    "/about/education",
	Handler: List,
	Guard: &guard.Guard{
		ParamStruct: &ListParams{},
	},
}

// ListParams represents the params accepted by the Add endpoint
type ListParams struct {
	paginator.HandlerParams
	Deleted  *bool  `from:"query" json:"deleted"`
	Orphans  *bool  `from:"query" json:"orphans"`
	Operator string `from:"query" json:"op" default:"and" enum:"and,or"`
}

// List is an endpoint used to list all Experience
func List(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*ListParams)
	paginator := params.Paginator()

	whereClause := makeListWhereClause(params, req.User())
	stmt := `SELECT edu.*, ` + organizations.JoinSQL("org") + `
	FROM about_education edu
	JOIN about_organizations org
	  ON org.id = edu.organization_id
	` + whereClause + `
	ORDER BY end_year DESC NULLS FIRST, start_year DESC
	OFFSET $1
	LIMIT $2`

	edus := ListEducation{}
	err := deps.DB.Select(&edus, stmt, paginator.Offset(), paginator.Limit())
	if err != nil {
		return err
	}

	if req.User().IsAdm() {
		return req.Response().Ok(edus.ExportPrivate())
	}
	return req.Response().Ok(edus.ExportPublic())
}

func makeListWhereClause(p *ListParams, u *auth.User) string {
	whereList := []string{}

	// Only the an admins can filter on deleted/orphans so regular
	// users only get non-orphans and non-deleted data
	if !u.IsAdm() {
		p.Operator = "and"
		whereList = append(whereList, "edu.deleted_at IS NULL")
		whereList = append(whereList, "org.deleted_at IS NULL")
	} else {
		if p.Orphans != nil {
			if *p.Orphans {
				whereList = append(whereList, "org.deleted_at IS NOT NULL")
			} else {
				whereList = append(whereList, "org.deleted_at IS NULL")
			}
		}

		if p.Deleted != nil {
			if *p.Deleted {
				whereList = append(whereList, "edu.deleted_at IS NOT NULL")
			} else {
				whereList = append(whereList, "edu.deleted_at IS NULL")
			}
		}
	}

	if len(whereList) > 0 {
		return "WHERE " + strings.Join(whereList, " "+p.Operator+" ")
	}
	return ""
}
