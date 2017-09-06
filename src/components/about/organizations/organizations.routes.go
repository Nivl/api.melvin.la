package organizations

import (
	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/gorilla/mux"
)

// Contains the index of all Endpoints
const (
	EndpointAdd = iota
	EndpointDelete
	EndpointUpdate
	EndpointList
	EndpointUploadLogo
)

// Endpoints is a list of endpoints for this components
var Endpoints = router.Endpoints{
	EndpointAdd:        addEndpoint,
	EndpointUpdate:     updateEndpoint,
	EndpointDelete:     deleteEndpoint,
	EndpointList:       listEndpoint,
	EndpointUploadLogo: uploadLogoEndpoint,
}

// SetRoutes is used to set all the organization routes
func SetRoutes(r *mux.Router, deps dependencies.Dependencies) {
	Endpoints.Activate(r, deps)
}
