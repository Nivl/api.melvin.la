package sessions

import (
	"github.com/Nivl/api.melvin.la/api/router"
	"github.com/gorilla/mux"
)

const (
	EndpointAdd = iota
	EndpointDelete
)

var Endpoints = router.Endpoints{
	EndpointAdd: {
		Verb:    "POST",
		Path:    "/",
		Handler: HandlerAdd,
		Params:  &HandlerAddParams{},
	},
	EndpointDelete: {
		Verb:    "DELETE",
		Path:    "/{token}",
		Handler: HandlerDelete,
		Params:  &HandlerDeleteParams{},
		Auth:    router.LoggedUser,
	},
}

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	Endpoints.Activate(r)
}
