package about

import (
	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/gorilla/mux"
	"github.com/melvin-laplanche/ml-api/src/components/about/education"
	"github.com/melvin-laplanche/ml-api/src/components/about/experience"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
)

// SetRoutes is used to set all the about routes
func SetRoutes(r *mux.Router, deps dependencies.Dependencies) {
	organizations.SetRoutes(r, deps)
	experience.SetRoutes(r, deps)
	education.SetRoutes(r, deps)
}
