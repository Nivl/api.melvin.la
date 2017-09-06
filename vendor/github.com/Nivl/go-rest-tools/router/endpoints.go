package router

import (
	"net/http"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/basicauth"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// Endpoints represents a list of endpoint
type Endpoints []*Endpoint

// Activate adds the endpoints to the router
func (endpoints Endpoints) Activate(router *mux.Router, apiDeps dependencies.Dependencies) {
	for _, endpoint := range endpoints {
		router.
			Methods(endpoint.Verb).
			Path(endpoint.Path).
			Handler(Handler(endpoint, apiDeps))
	}
}

// Handler makes it possible to use a RouteHandler where a http.Handler is required
func Handler(e *Endpoint, apiDeps dependencies.Dependencies) http.Handler {
	HTTPHandler := func(resWriter http.ResponseWriter, req *http.Request) {
		// storageErr will be checked later on. Since the storage is not
		// needed to throw an error, we first init the request, then we will use
		// that request to return (and log) the error
		fileStorage, storageErr := apiDeps.FileStorage(req.Context())

		handlerDeps := &Dependencies{
			DB:      apiDeps.DB(),
			Storage: fileStorage,
			Mailer:  apiDeps.Mailer(),
		}
		request := &Request{
			id:     uuid.NewV4().String()[:8],
			http:   req,
			res:    NewResponse(resWriter, handlerDeps),
			logger: apiDeps.Logger(),
		}
		defer request.handlePanic()

		// We set some response data
		request.res.Header().Set("X-Request-Id", request.id)

		// if we failed getting the dependencies, we return a 500
		if storageErr != nil {
			request.res.Error(storageErr, request)
			return
		}

		// We fetch the user session if a token is provided
		headers, found := req.Header["Authorization"]
		if found {
			userID, sessionID, err := basicauth.ParseAuthHeader(headers, "basic", "")
			if err != nil {
				request.res.Error(apierror.NewBadRequest("Authorization", "invalid format"), request)
			}
			session := &auth.Session{ID: sessionID, UserID: userID}

			if session.ID != "" && session.UserID != "" {
				exists, err := session.Exists(handlerDeps.DB)
				if err != nil {
					request.res.Error(err, request)
					return
				}
				if !exists {
					request.res.Error(apierror.NewNotFoundField("Authorization", "session not found"), request)
					return
				}
				request.session = session
				// we get the user and make sure it (still) exists
				request.user, err = auth.GetUserByID(handlerDeps.DB, session.UserID)
				if err != nil {
					if apierror.IsNotFound(err) {
						err = apierror.NewNotFoundField("Authorization", "session not found")
					}
					request.res.Error(err, request)
					return
				}
			}
		}

		// Make sure the user has access to the handler
		if allowed, err := e.Guard.HasAccess(request.user); !allowed {
			request.res.Error(err, request)
			return
		}

		// We Parse the request params
		if e.Guard != nil && e.Guard.ParamStruct != nil {
			// Get the list of all http params provided by the client
			sources, err := request.httpParamsBySource()
			if err != nil {
				request.res.Error(err, request)
				return
			}

			request.params, err = e.Guard.ParseParams(sources, request.http)
			if err != nil {
				request.res.Error(err, request)
				return
			}
		}

		// Execute the actual route handler
		err := e.Handler(request, handlerDeps)
		if err != nil {
			request.res.Error(err, request)
		}
	}

	return http.HandlerFunc(HTTPHandler)
}
