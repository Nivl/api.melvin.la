package sessions

import (
	"github.com/gorilla/mux"
	"github.com/melvin-laplanche/ml-api/src/router"
)

// Contains the index of all Endpoints
const (
	EndpointAdd = iota
	EndpointDelete
)

// Endpoints is a list of endpoints for this components
var Endpoints = router.Endpoints{
	EndpointAdd: {
		Verb:    "POST",
		Path:    "/sessions",
		Handler: HandlerAdd,
		Params:  &HandlerAddParams{},
	},
	EndpointDelete: {
		Verb:    "DELETE",
		Path:    "/sessions/{token}",
		Handler: HandlerDelete,
		Params:  &HandlerDeleteParams{},
		Auth:    router.LoggedUser,
	},
}

// SetRoutes is used to set all the routes of the article
func SetRoutes(baseURI string, r *mux.Router) {
	Endpoints.Activate(baseURI, r)
}
