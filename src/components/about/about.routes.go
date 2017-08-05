package about

import "github.com/gorilla/mux"
import "github.com/melvin-laplanche/ml-api/src/components/about/organizations"

func SetRoutes(r *mux.Router) {
	organizations.SetRoutes(r)
}
