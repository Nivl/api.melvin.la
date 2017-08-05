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
func (endpoints Endpoints) Activate(router *mux.Router) {
	for _, endpoint := range endpoints {
		router.
			Methods(endpoint.Verb).
			Path(endpoint.Path).
			Handler(Handler(endpoint))
	}
}

// Handler makes it possible to use a RouteHandler where a http.Handler is required
func Handler(e *Endpoint) http.Handler {
	HTTPHandler := func(resWriter http.ResponseWriter, req *http.Request) {
		deps, depsErr := NewDefaultDependenciesWithContext(req.Context())
		// we need deps for the response, so if it fails we get only the main
		// deps that we will use for the response, then we'll use that response
		// to return a 500
		if depsErr != nil {
			deps = NewNoFailersDependencies()
		}
		request := &Request{
			id:     uuid.NewV4().String()[:8],
			http:   req,
			res:    NewResponse(resWriter, deps),
			logger: dependencies.NewLogger(),
		}
		defer request.handlePanic()

		// We set some response data
		request.res.Header().Set("X-Request-Id", request.id)

		// if we failed getting the dependencies, we return a 500
		if depsErr != nil {
			request.res.Error(depsErr, request)
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
				exists, err := session.Exists(deps.DB)
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
				request.user, err = auth.GetUserByID(deps.DB, session.UserID)
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
		err := e.Handler(request, deps)
		if err != nil {
			request.res.Error(err, request)
		}
	}

	return http.HandlerFunc(HTTPHandler)
}
