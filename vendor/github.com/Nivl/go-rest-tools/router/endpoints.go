package router

import (
	"net/http"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/logger"
	"github.com/Nivl/go-rest-tools/network/http/basicauth"
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// Endpoints represents a list of endpoint
type Endpoints []*Endpoint

// Activate adds the endpoints to the router
func (endpoints Endpoints) Activate(router *mux.Router) {
	deps := NewDefaultDependencies()

	for _, endpoint := range endpoints {
		router.
			Methods(endpoint.Verb).
			Path(endpoint.Path).
			Handler(Handler(endpoint, deps))
	}
}

// Handler makes it possible to use a RouteHandler where a http.Handler is required
func Handler(e *Endpoint, deps *Dependencies) http.Handler {
	HTTPHandler := func(resWriter http.ResponseWriter, req *http.Request) {
		request := &Request{
			id:   uuid.NewV4().String()[:8],
			http: req,
			res:  NewResponse(resWriter, deps),
			deps: deps,
		}
		defer request.handlePanic()

		if dependencies.Logentries != nil {
			request.logger = logger.NewLogEntries(dependencies.Logentries)
		} else {
			request.logger = logger.NewBasicLogger()
		}

		// We set some response data
		request.res.Header().Set("X-Request-Id", request.id)

		// We Parse the request params
		if e.Guard != nil && e.Guard.ParamStruct != nil {
			// Get the list of all http params provided by the client
			sources, err := request.httpParamsBySource()
			if err != nil {
				request.res.Error(err, request)
				return
			}

			request.params, err = e.Guard.ParseParams(sources)
			if err != nil {
				request.res.Error(err, request)
				return
			}
		}

		// We check the auth
		headers, found := req.Header["Authorization"]
		if found {
			userID, sessionID, _ := basicauth.ParseAuthHeader(headers, "basic", "")
			session := &auth.Session{ID: sessionID, UserID: userID}

			if session.ID != "" && session.UserID != "" {
				exists, err := session.Exists(deps.DB)
				if err != nil {
					request.res.Error(err, request)
					return
				}
				if !exists {
					request.res.Error(httperr.NewBadRequest("invalid auth data"), request)
					return
				}
				// we get the user and make sure it (still) exists
				request.user, err = auth.GetUser(deps.DB, session.UserID)
				if err != nil {
					request.res.Error(err, request)
					return
				}
				if request.user == nil {
					request.res.Error(httperr.NewBadRequest("user not found"), request)
					return
				}
			}
		}

		if allowed, err := e.Guard.HasAccess(request.user); !allowed {
			request.res.Error(err, request)
			return
		}

		// Execute the actual route handler
		err := e.Handler(request, deps)
		if err != nil {
			request.res.Error(err, request)
		}
	}

	return http.HandlerFunc(HTTPHandler)
}
