package articles

import (
	"github.com/Nivl/go-rest-tools/router"
	"github.com/gorilla/mux"
)

// Indexes of all different endpoints
const (
	EndpointAdd = iota
	EndpointGet
	EndpointSearch
	EndpointUpdate
	EndpointDelete
	EndpointVersionAdd
	EndpointVersionList
	EndpointVersionUpdate
	EndpointVersionDelete
)

// Endpoints contains the list of endpoints for this component
var Endpoints = router.Endpoints{
	EndpointAdd: {
		Verb:    "POST",
		Path:    "/blog/articles",
		Handler: Add,
		Auth:    router.AdminAccess,
		Params:  &AddParams{},
	},
	EndpointGet: {
		Verb:    "GET",
		Path:    "/blog/articles/{id}",
		Handler: Get,
		Auth:    nil,
		Params:  &GetParams{},
	},
	EndpointSearch: {
		Verb:    "GET",
		Path:    "/blog/articles",
		Handler: Search,
		Auth:    nil,
		Params:  &SearchParams{},
	},
	EndpointUpdate: {
		Verb:    "PATCH",
		Path:    "/blog/articles/{id}",
		Handler: Update,
		Auth:    router.AdminAccess,
		Params:  &UpdateParams{},
	},
	EndpointDelete: {
		Verb:    "DELETE",
		Path:    "/blog/articles/{id}",
		Handler: Delete,
		Auth:    router.AdminAccess,
		Params:  &DeleteParams{},
	},
	EndpointVersionAdd: {
		Verb:    "POST",
		Path:    "/blog/articles/{article_id}/versions",
		Handler: AddVersion,
		Auth:    router.AdminAccess,
		Params:  &AddVersionParams{},
	},
	EndpointVersionList: {
		Verb:    "GET",
		Path:    "/blog/articles/{article_id}/versions",
		Handler: ListVersion,
		Auth:    router.AdminAccess,
		Params:  &ListVersionParams{},
	},
	EndpointVersionUpdate: {
		Verb:    "PATCH",
		Path:    "/blog/articles/{article_id}/versions/{id}",
		Handler: UpdateVersion,
		Auth:    router.AdminAccess,
		Params:  &UpdateVersionParams{},
	},
	EndpointVersionDelete: {
		Verb:    "DELETE",
		Path:    "/blog/articles/{article_id}/versions/{id}",
		Handler: DeleteVersion,
		Auth:    router.AdminAccess,
		Params:  &DeleteVersionParams{},
	},
}

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	Endpoints.Activate(r)
}
