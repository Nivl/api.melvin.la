package experience

import (
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
)

var getEndpoint = &router.Endpoint{
	Verb:    "GET",
	Path:    "/about/experience/{id}",
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

	var exp *Experience
	var err error

	// admins
	if req.User().IsAdm() {
		exp, err = GetAnyByID(deps.DB, params.ID)
	} else {
		exp, err = GetByID(deps.DB, params.ID)
	}

	if err != nil {
		return err
	}

	var pld *Payload
	if req.User().IsAdm() {
		pld = exp.ExportPrivate()
	} else {
		pld = exp.ExportPublic()
	}

	return req.Response().Ok(pld)
}
