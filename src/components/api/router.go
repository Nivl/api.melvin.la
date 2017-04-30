package api

import (
	"fmt"
	"net/http"

	"github.com/Nivl/go-rest-tools/network/http/httpres"
	"github.com/gorilla/mux"
	"github.com/melvin-laplanche/ml-api/src/components/blog"
	"github.com/melvin-laplanche/ml-api/src/components/sessions"
	"github.com/melvin-laplanche/ml-api/src/components/users"
)

func notFound(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := fmt.Sprintf(`{"error":"%s"}`, http.StatusText(http.StatusNotFound))
	httpres.ErrorJSON(w, err, http.StatusNotFound)
}

// GetRouter return the api router with all the routes
func GetRouter() *mux.Router {
	r := mux.NewRouter()
	blog.SetRoutes(r)
	users.SetRoutes(r)
	sessions.SetRoutes(r)
	r.NotFoundHandler = http.HandlerFunc(notFound)

	return r
}
