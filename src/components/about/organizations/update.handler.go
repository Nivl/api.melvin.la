package organizations

import (
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-types/datetime"
)

var updateEndpoint = &router.Endpoint{
	Verb:    "PATCH",
	Path:    "/about/organizations/{id}",
	Handler: Update,
	Guard: &guard.Guard{
		ParamStruct: &UpdateParams{},
		Auth:        guard.AdminAccess,
	},
}

// UpdateParams represents the request params accepted by the Update handler
type UpdateParams struct {
	ID        string  `from:"url" json:"id" params:"required,uuid"`
	Name      *string `from:"form" json:"name" params:"trim,noempty"`
	ShortName *string `from:"form" json:"short_name" params:"trim"`
	Website   *string `from:"form" json:"website" params:"url,trim"`
	InTrash   *bool   `from:"form" json:"in_trash"`
}

// Update represent an API handler to update an organization
func Update(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*UpdateParams)

	org, err := GetAnyByID(deps.DB, params.ID)
	if err != nil {
		return err
	}

	if params.Name != nil {
		org.Name = *params.Name
	}
	if params.ShortName != nil {
		org.ShortName = params.ShortName
	}
	if params.Website != nil {
		org.Website = params.Website
	}
	if params.InTrash != nil {
		if *params.InTrash && org.DeletedAt == nil {
			org.DeletedAt = datetime.Now()
		} else {
			org.DeletedAt = nil
		}
	}

	if err := org.Update(deps.DB); err != nil {
		return err
	}

	req.Response().Ok(org.ExportPrivate())
	return nil
}
