package api

import (
	"github.com/melvin-laplanche/ml-api/src/components/blog"
	"github.com/melvin-laplanche/ml-api/src/components/sessions"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	r := mux.NewRouter()
	r.Host("api.melvin.la")
	r.Host("api.melvin.loc")
	blog.SetRoutes(r.PathPrefix("/blog").Subrouter())
	users.SetRoutes(r.PathPrefix("/users").Subrouter())
	sessions.SetRoutes(r.PathPrefix("/sessions").Subrouter())
	//router.NotFoundHandler = http.HandlerFunc(noRoutes)

	return r
}
