package about

import (
	"github.com/gorilla/mux"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
)

func SetRoutes(r *mux.Router) {
	organizations.SetRoutes(r)
	experience.SetRoutes(r)
}
