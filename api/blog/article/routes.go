package article

import (
	"github.com/Nivl/api.melvin.la/api/blog/article/articlehandlers"
	"github.com/Nivl/api.melvin.la/api/router"
	"github.com/gorilla/mux"
)

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	routes := router.Endpoints{
		{
			Verb:    "GET",
			Path:    "/",
			Handler: articlehandlers.List,
			Auth:    nil,
		},
		{
			Verb:    "GET",
			Path:    "/{id}",
			Handler: articlehandlers.GetOne,
			Auth:    nil,
		},
		{
			Verb:    "POST",
			Path:    "/",
			Handler: articlehandlers.Add,
			Auth:    nil,
		},
		{
			Verb:    "PATCH",
			Path:    "/{id}",
			Handler: articlehandlers.Update,
			Auth:    nil,
		},
	}

	routes.Activate(r)
}
