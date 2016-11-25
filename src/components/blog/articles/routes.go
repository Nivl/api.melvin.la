package articles

import (
	"github.com/gorilla/mux"
	"github.com/melvin-laplanche/ml-api/src/router"
)

const (
	EndpointAdd = iota
)

var Endpoints = router.Endpoints{
	EndpointAdd: {
		Verb:    "POST",
		Path:    "/",
		Handler: HandlerAdd,
		Auth:    router.LoggedUser,
		Params:  &HandlerAddParams{},
	},
}

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	Endpoints.Activate(r)
}
