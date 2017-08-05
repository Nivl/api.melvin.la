package api

import (
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/gorilla/mux"
	"github.com/melvin-laplanche/ml-api/src/components/about"
	"github.com/melvin-laplanche/ml-api/src/components/sessions"
	"github.com/melvin-laplanche/ml-api/src/components/users"
)

var notFoundEndpoint = &router.Endpoint{
	Handler: func(req router.HTTPRequest, deps *router.Dependencies) error {
		return apierror.NewNotFound()
	},
}

// GetRouter return the api router with all the routes
func GetRouter() *mux.Router {
	r := mux.NewRouter()
	users.SetRoutes(r)
	sessions.SetRoutes(r)
	about.SetRoutes(r)
	r.NotFoundHandler = router.Handler(notFoundEndpoint)
	return r
}
