package blog

import (
	"github.com/Nivl/api.melvin.la/api/blog/article"
	"github.com/gorilla/mux"
)

// SetRoutes is used to set all the routes of the blog
func SetRoutes(r *mux.Router) {
	article.SetRoutes(r.PathPrefix("/articles").Subrouter())
}
