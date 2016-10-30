package api

import (
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/components/blog"
	"github.com/Nivl/api.melvin.la/api/components/users"
	"github.com/gorilla/mux"
)

func EnsureIndexes() {
	auth.EnsureIndexes()
	blog.EnsureIndexes()
}

func GetRouter() *mux.Router {
	r := mux.NewRouter()
	r.Host("api.melvin.la")
	r.Host("api.melvin.loc")
	blog.SetRoutes(r.PathPrefix("/blog").Subrouter())
	users.SetRoutes(r.PathPrefix("/users").Subrouter())
	//router.NotFoundHandler = http.HandlerFunc(noRoutes)

	return r
}
