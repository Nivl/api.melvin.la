package router

import (
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/auth"
	uuid "github.com/satori/go.uuid"
)

// Endpoints represents a list of endpoint
type Endpoints []*Endpoint

// Activate adds the endpoints to the router
func (endpoints Endpoints) Activate(basePath string, router *mux.Router) {
	for _, endpoint := range endpoints {
		// Get the full path
		fullPath := basePath + endpoint.Path
		if endpoint.Prefix != "" {
			fullPath = endpoint.Prefix + fullPath
		}

		router.
			Methods(endpoint.Verb).
			Path(fullPath).
			Handler(Handler(endpoint))
	}
}

// Handler makes it possible to use a RouteHandler where a http.Handler is required
func Handler(e *Endpoint) http.Handler {
	HTTPHandler := func(resWriter http.ResponseWriter, req *http.Request) {
		request := &Request{
			ID:       uuid.NewV4().String()[:8],
			Request:  req,
			Response: resWriter,
		}
		defer request.handlePanic()

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
		session := &auth.Session{ID: req.Header.Get("X-Session-Token"), UserID: req.Header.Get("X-User-Id")}
		if session.ID != "" && session.UserID != "" {
			exists, err := session.Exists()
			if err != nil {
				request.Error(err)
				return
			}
			if !exists {
				request.Error(apierror.NewBadRequest("invalid auth data"))
				return
			}
			// we get the user and make sure it (still) exists
			request.User, err = auth.GetUser(session.UserID)
			if err != nil {
				request.Error(err)
				return
			}
			if request.User == nil {
				request.Error(apierror.NewBadRequest("user not found"))
				return
			}
		}

		// We set some response data
		request.Response.Header().Set("X-Request-Id", request.ID)

		accessGranted := e.Auth == nil || e.Auth(request)
		if !accessGranted {
			request.Error(apierror.NewUnauthorized())
			return
		}

		// Execute the actual route handler
		err := e.Handler(request)
		if err != nil {
			request.Error(err)
		}
	}

	return http.HandlerFunc(HTTPHandler)
}
