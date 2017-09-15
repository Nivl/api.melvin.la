package users

import (
	"net/http"

	"github.com/Nivl/go-rest-tools/router"
)

var getFeaturedEndpoint = &router.Endpoint{
	Verb:    http.MethodGet,
	Path:    "/users/featured",
	Handler: GetFeatured,
}

// GetFeatured represent an API handler to get a user
func GetFeatured(req router.HTTPRequest, deps *router.Dependencies) error {
	profile, err := GetFeaturedProfile(deps.DB)
	if err != nil {
		return err
	}
	return req.Response().Ok(profile.ExportPublic())
}
