package users

import (
	"github.com/Nivl/go-rest-tools/router"
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
		Auth:    nil,
		Handler: Add,
		Params:  &AddParams{},
	},
	EndpointUpdate: {
		Verb:    "PATCH",
		Path:    "/users/{id}",
		Auth:    router.LoggedUserAccess,
		Handler: Update,
		Params:  &UpdateParams{},
	},
	EndpointDelete: {
		Verb:    "DELETE",
		Path:    "/users/{id}",
		Auth:    router.LoggedUserAccess,
		Handler: Delete,
		Params:  &DeleteParams{},
	},
	EndpointGet: {
		Verb:    "GET",
		Path:    "/users/{id}",
		Auth:    nil,
		Handler: Get,
		Params:  &GetParams{},
	},
}

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	Endpoints.Activate(r)
}
