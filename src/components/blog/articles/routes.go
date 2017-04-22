package articles

import (
	"github.com/Nivl/go-rest-tools/router"
	"github.com/gorilla/mux"
)

// Indexes of all different endpoints
const (
	EndpointAdd = iota
	EndpointGet
	// EndpointList
	EndpointUpdate
	EndpointDelete
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
	// EndpointList: {
	// 	Verb: "GET",
	// 	Path: "/articles",
	// 	//Handler: HandlerList,
	// 	Auth: nil,
	// },
	// EndpointUpdateDraft: {
	// 	Verb:    "PATCH",
	// 	Path:    "/articles/{id}/draft",
	// 	Handler: HandlerUpdateDraft,
	// 	Auth:    router.LoggedUser,
	// 	Params:  &HandlerUpdateDraftParams{},
	// },
	// EndpointUserList: {
	// 	Verb:    "GET",
	// 	Prefix:  "/users/{user_id}",
	// 	Path:    "/articles",
	// 	Handler: HandlerListForUser,
	// 	Auth:    nil,
	// 	Params:  &HandlerListForUserParams{},
	// },
}

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	Endpoints.Activate(r)
}
