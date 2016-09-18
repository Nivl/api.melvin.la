package article

import (
	"github.com/Nivl/api.melvin.la/src/router"
	"github.com/gorilla/mux"
)

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	routes := router.Endpoints{
		{
			Verb:    "GET",
			Path:    "/",
			Handler: List,
			Auth:    nil,
		},
		{
			Verb:    "GET",
			Path:    "/{id}",
			Handler: GetOne,
			Auth:    nil,
		},
		{
			Verb:    "POST",
			Path:    "/",
			Handler: Add,
			Auth:    nil,
		},
		{
			Verb:    "PATCH",
			Path:    "/{id}",
			Handler: Update,
			Auth:    nil,
		},
	}

	routes.Activate(r)
}
