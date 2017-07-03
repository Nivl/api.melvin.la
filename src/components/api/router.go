package api

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/gorilla/mux"
	"github.com/melvin-laplanche/ml-api/src/components/users"
)

var notFoundEndpoint = &router.Endpoint{
	Handler: func(req router.HTTPRequest, deps *router.Dependencies) error {
		return httperr.NewNotFound()
	},
}

// GetRouter return the api router with all the routes
func GetRouter() *mux.Router {
	r := mux.NewRouter()
	// blog.SetRoutes(r)
	users.SetRoutes(r)
	// sessions.SetRoutes(r)
	r.NotFoundHandler = router.Handler(notFoundEndpoint, router.NewDefaultDependencies())

	return r
}
