package experience

import (
	"github.com/Nivl/go-rest-tools/router"
	"github.com/gorilla/mux"
)

// Contains the index of all Endpoints
const (
	EndpointAdd = iota
	EndpointGet
)

// Endpoints is a list of endpoints for this components
var Endpoints = router.Endpoints{
	EndpointAdd: addEndpoint,
	EndpointGet: getEndpoint,
}

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	Endpoints.Activate(r)
}
