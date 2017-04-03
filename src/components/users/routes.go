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
		Handler: HandlerAdd,
		Params:  &HandlerAddParams{},
	},
	EndpointUpdate: {
		Verb:    "PATCH",
		Path:    "/users/{id}",
		Auth:    router.LoggedUserAccess,
		Handler: HandlerUpdate,
		Params:  &HandlerUpdateParams{},
	},
	EndpointDelete: {
		Verb:    "DELETE",
		Path:    "/users/{id}",
		Auth:    router.LoggedUserAccess,
		Handler: HandlerDelete,
		Params:  &HandlerDeleteParams{},
	},
	EndpointGet: {
		Verb:    "GET",
		Path:    "/users/{id}",
		Auth:    nil,
		Handler: HandlerGet,
		Params:  &HandlerGetParams{},
	},
}

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	Endpoints.Activate(r)
}
