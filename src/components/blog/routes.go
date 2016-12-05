package blog

import (
	"github.com/gorilla/mux"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles"
)

// SetRoutes is used to set all the routes of the blog
func SetRoutes(baseURI string, r *mux.Router) {
	baseURI += "/blog"
	articles.SetRoutes(baseURI, r)
}
