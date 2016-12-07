package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/melvin-laplanche/ml-api/src/components/blog"
	"github.com/melvin-laplanche/ml-api/src/components/sessions"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/melvin-laplanche/ml-api/src/mlhttp"
)

// AllowedOrigins is a list containing all origins allowed to hit the API
var AllowedOrigins = []string{
	"http://www.melvin.la",
	"http://orchid.melvin.la",
}

// AllowedMethods is a list containing all HTTP verb accepted by the API
var AllowedMethods = []string{
	"GET", "POST", "PATCH", "DELETE",
}

// AllowedHeaders is a list custom headers accepted by the API
var AllowedHeaders = []string{
	"X-Session-Token", "X-User-Id",
}

func notFound(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := fmt.Sprintf(`{"error":"%s"}`, http.StatusText(http.StatusNotFound))
	mlhttp.ErrorJSON(w, err, http.StatusNotFound)
}

// GetRouter return the api router with all the routes
func GetRouter() *mux.Router {
	baseURI := ""
	r := mux.NewRouter()
	blog.SetRoutes(baseURI, r)
	users.SetRoutes(baseURI, r)
	sessions.SetRoutes(baseURI, r)
	r.NotFoundHandler = http.HandlerFunc(notFound)

	return r
}
