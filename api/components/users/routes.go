package users

import (
	"github.com/Nivl/api.melvin.la/api/router"
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
		Path:    "/",
		Auth:    nil,
		Handler: HandlerAdd,
		Params:  &HandlerAddParams{},
	},
	EndpointUpdate: {
		Verb:    "PATCH",
		Path:    "/{id}",
		Auth:    router.LoggedUser,
		Handler: HandlerUpdate,
		Params:  &HandlerUpdateParams{},
	},
	EndpointDelete: {
		Verb:    "DELETE",
		Path:    "/{id}",
		Auth:    router.LoggedUser,
		Handler: HandlerDelete,
		Params:  &HandlerDeleteParams{},
	},
	EndpointGet: {
		Verb:    "GET",
		Path:    "/{id}",
		Auth:    nil,
		Handler: HandlerGet,
		Params:  &HandlerGetParams{},
	},
}

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	Endpoints.Activate(r)
}
