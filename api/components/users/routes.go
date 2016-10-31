package users

import (
	"github.com/Nivl/api.melvin.la/api/router"
	"github.com/gorilla/mux"
)

const (
	EndpointAdd = iota
	EndpointUpdate
	EndpointDelete
)

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
}

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	Endpoints.Activate(r)
}
