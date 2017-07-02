package users

import (
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/gorilla/mux"
)

// Contains the index of all Endpoints
const (
	EndpointAdd = iota
	EndpointUpdate
	EndpointDelete
	EndpointGet
)

// Endpoints is a list of endpoints for this components
var Endpoints = router.Endpoints{
	EndpointAdd: {
		Verb:    "POST",
		Path:    "/users",
		Handler: Add,
		Guard: &guard.Guard{
			ParamStruct: &AddParams{},
		},
	},
	EndpointUpdate: {
		Verb:    "PATCH",
		Path:    "/users/{id}",
		Handler: Update,
		Guard: &guard.Guard{
			ParamStruct: &UpdateParams{},
			Auth:        guard.LoggedUserAccess,
		},
	},
	EndpointDelete: {
		Verb:    "DELETE",
		Path:    "/users/{id}",
		Handler: Delete,
		Guard: &guard.Guard{
			ParamStruct: &DeleteParams{},
			Auth:        guard.LoggedUserAccess,
		},
	},
	EndpointGet: {
		Verb:    "GET",
		Path:    "/users/{id}",
		Handler: Get,
		Guard: &guard.Guard{
			ParamStruct: &GetParams{},
		},
	},
}

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	Endpoints.Activate(r)
}
