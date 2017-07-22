package organizations

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
)

var deleteEndpoint = &router.Endpoint{
	Verb:    "DELETE",
	Path:    "/about/organizations/{id}",
	Handler: Delete,
	Guard: &guard.Guard{
		ParamStruct: &DeleteParams{},
		Auth:        guard.AdminAccess,
	},
}

// DeleteParams represent the request params accepted by the Delete handler
type DeleteParams struct {
	ID string `from:"url" json:"id" params:"required,uuid"`
}

// Delete represent an API handler to remove a session
func Delete(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*DeleteParams)

	org, err := GetByID(deps.DB, params.ID)
	if err != nil {
		return err
	}

	if org.IsZero() {
		return httperr.NewNotFound()
	}

	if err := org.Delete(deps.DB); err != nil {
		return err
	}

	req.Response().NoContent()
	return nil
}
