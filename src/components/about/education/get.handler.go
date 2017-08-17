package education

import (
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
)

var getEndpoint = &router.Endpoint{
	Verb:    "GET",
	Path:    "/about/education/{id}",
	Handler: Get,
	Guard: &guard.Guard{
		ParamStruct: &GetParams{},
	},
}

// GetParams represents the params accepted by the Add endpoint
type GetParams struct {
	ID string `from:"url" json:"id" params:"required,uuid"`
}

// Get is an endpoint used to get an Organization
func Get(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*GetParams)

	var edu *Education
	var err error

	// admins
	if req.User().IsAdm() {
		edu, err = GetAnyByID(deps.DB, params.ID)
	} else {
		edu, err = GetByID(deps.DB, params.ID)
	}

	if err != nil {
		return err
	}

	var pld *Payload
	if req.User().IsAdm() {
		pld = edu.ExportPrivate()
	} else {
		pld = edu.ExportPublic()
	}

	return req.Response().Ok(pld)
}
