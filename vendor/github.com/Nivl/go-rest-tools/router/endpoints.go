package router

import (
	"net/http"
	"reflect"

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
			ID:       uuid.NewV4().String()[:8],
			Request:  req,
			Response: resWriter,
			deps:     deps,
		}
		defer request.handlePanic()

		if dependencies.Logentries != nil {
			request.Logger = logger.NewLogEntries(dependencies.Logentries)
		} else {
			request.Logger = logger.NewBasicLogger()
		}

		// We set some response data
		request.Response.Header().Set("X-Request-Id", request.ID)

		// We Parse the request params
		if e.Params != nil {
			// We give request.Params the same type as e.Params
			request.Params = reflect.New(reflect.TypeOf(e.Params).Elem()).Interface()
			if err := request.ParseParams(); err != nil {
				request.Error(err)
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
					request.Error(err)
					return
				}
				if !exists {
					request.Error(httperr.NewBadRequest("invalid auth data"))
					return
				}
				// we get the user and make sure it (still) exists
				request.User, err = auth.GetUser(deps.DB, session.UserID)
				if err != nil {
					request.Error(err)
					return
				}
				if request.User == nil {
					request.Error(httperr.NewBadRequest("user not found"))
					return
				}
			}
		}

		accessGranted := e.Auth == nil || e.Auth(request)
		if !accessGranted {
			request.Error(httperr.NewUnauthorized())
			return
		}

		// Execute the actual route handler
		err := e.Handler(request, deps)
		if err != nil {
			request.Error(err)
		}
	}

	return http.HandlerFunc(HTTPHandler)
}
