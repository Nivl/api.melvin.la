package blog

import (
	"github.com/Nivl/api.melvin.la/api/components/blog/articles"
	"github.com/gorilla/mux"
)

// SetRoutes is used to set all the routes of the blog
func SetRoutes(r *mux.Router) {
	articles.SetRoutes(r.PathPrefix("/articles").Subrouter())
}
