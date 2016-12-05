package users

import (
	"github.com/gorilla/mux"
	"github.com/melvin-laplanche/ml-api/src/router"
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
		Auth:    router.LoggedUser,
		Handler: HandlerUpdate,
		Params:  &HandlerUpdateParams{},
	},
	EndpointDelete: {
		Verb:    "DELETE",
		Path:    "/users/{id}",
		Auth:    router.LoggedUser,
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
func SetRoutes(baseURI string, r *mux.Router) {
	Endpoints.Activate(baseURI, r)
}
