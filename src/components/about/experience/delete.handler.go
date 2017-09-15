package experience

import (
	"net/http"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
)

var deleteEndpoint = &router.Endpoint{
	Verb:    http.MethodDelete,
	Path:    "/about/experience/{id}",
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

	exp, err := GetAnyByID(deps.DB, params.ID)
	if err != nil {
		return err
	}

	if err := exp.Delete(deps.DB); err != nil {
		return err
	}

	req.Response().NoContent()
	return nil
}
