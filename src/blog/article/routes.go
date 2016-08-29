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
			Handler: ArticleList,
			Auth:    nil,
		},
		{
			Verb:    "GET",
			Path:    "/{id}",
			Handler: ArticleGet,
			Auth:    nil,
		},
		{
			Verb:             "POST",
			Path:             "/",
			Handler:          ArticleAdd,
			Auth:             nil,
			JSONBodyTemplate: &Article{},
		},
		{
			Verb:             "PATCH",
			Path:             "/{id}",
			Handler:          ArticleUpdate,
			Auth:             nil,
			JSONBodyTemplate: &Article{},
		},
	}

	routes.Activate(r)
}
