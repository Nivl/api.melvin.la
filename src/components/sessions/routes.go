package sessions

import (
	"github.com/Nivl/go-rest-tools/router"
	"github.com/gorilla/mux"
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
		Handler: Add,
		Params:  &AddParams{},
	},
	EndpointDelete: {
		Verb:    "DELETE",
		Path:    "/sessions/{token}",
		Handler: Delete,
		Params:  &DeleteParams{},
		Auth:    router.LoggedUserAccess,
	},
}

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	Endpoints.Activate(r)
}
