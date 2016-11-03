package articles

import (
	"github.com/Nivl/api.melvin.la/api/router"
	"github.com/gorilla/mux"
)

const (
	EndpointList = iota
	EndpointGet
	EndpointAdd
	EndpointUpdate
)

var Endpoints = router.Endpoints{
	EndpointList: {
		Verb:    "GET",
		Path:    "/",
		Handler: HandlerList,
		Auth:    nil,
	},
	EndpointGet: {
		Verb:    "GET",
		Path:    "/{id}",
		Handler: HandlerGet,
		Auth:    nil,
	},
	EndpointAdd: {
		Verb:    "POST",
		Path:    "/",
		Handler: HandlerAdd,
		Auth:    router.LoggedUser,
		Params:  &HandlerAddParams{},
	},
	EndpointUpdate: {
		Verb:    "PATCH",
		Path:    "/{id}",
		Handler: HandlerUpdate,
		Auth:    router.LoggedUser,
	},
}

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	Endpoints.Activate(r)
}
