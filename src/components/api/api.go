package api

import (
	"github.com/gorilla/mux"
	"github.com/melvin-laplanche/ml-api/src/components/blog"
	"github.com/melvin-laplanche/ml-api/src/components/sessions"
	"github.com/melvin-laplanche/ml-api/src/components/users"
)

// GetRouter return the api router with all the routes
func GetRouter() *mux.Router {
	baseURI := ""
	r := mux.NewRouter()
	blog.SetRoutes(baseURI, r)
	users.SetRoutes(baseURI, r)
	sessions.SetRoutes(baseURI, r)
	//router.NotFoundHandler = http.HandlerFunc(noRoutes)

	return r
}
