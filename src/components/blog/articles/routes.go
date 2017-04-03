package articles

import (
	"github.com/Nivl/go-rest-tools/router"
	"github.com/gorilla/mux"
)

// Indexes of all different endpoints
const (
	EndpointAdd = iota
	EndpointList
	EndpointGet
	EndpointUpdate
	EndpointUpdateDraft
	EndpointDelete
	EndpointDeleteDraft
	EndpointUserList
)

// Endpoints contains the list of endpoints for this component
var Endpoints = router.Endpoints{
// EndpointAdd: {
// 	Verb:    "POST",
// 	Path:    "/articles",
// 	Handler: HandlerAdd,
// 	Auth:    router.LoggedUser,
// 	Params:  &HandlerAddParams{},
// },
// EndpointList: {
// 	Verb:    "GET",
// 	Path:    "/articles",
// 	Handler: HandlerList,
// 	Auth:    nil,
// },
// EndpointGet: {
// 	Verb:    "GET",
// 	Path:    "/articles/{id}",
// 	Handler: HandlerGet,
// 	Auth:    nil,
// 	Params:  &HandlerGetParams{},
// },
// EndpointUpdate: {
// 	Verb:    "PATCH",
// 	Path:    "/articles/{id}",
// 	Handler: HandlerUpdate,
// 	Auth:    router.LoggedUser,
// 	Params:  &HandlerUpdateParams{},
// },
// EndpointDelete: {
// 	Verb:    "DELETE",
// 	Path:    "/articles/{id}",
// 	Handler: HandlerDelete,
// 	Auth:    router.LoggedUser,
// 	Params:  &HandlerDeleteParams{},
// },
// EndpointUpdateDraft: {
// 	Verb:    "PATCH",
// 	Path:    "/articles/{id}/draft",
// 	Handler: HandlerUpdateDraft,
// 	Auth:    router.LoggedUser,
// 	Params:  &HandlerUpdateDraftParams{},
// },
// EndpointDeleteDraft: {
// 	Verb:    "DELETE",
// 	Path:    "/articles/{id}/draft",
// 	Handler: HandlerDeleteDraft,
// 	Auth:    router.LoggedUser,
// 	Params:  &HandlerDeleteDraftParams{},
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
