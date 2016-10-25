package router

import (
	"net/http"
	"reflect"

	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type Endpoints []*Endpoint

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
		sess := auth.NewSessionFromStrings(req.Header.Get("X-Session-Id"), req.Header.Get("X-User-Id"))
		if sess.ID != "" && sess.UserID != "" {
			exists, err := sess.Exists()
			if err != nil {
				request.Error(err)
				return
			}
			if !exists {
				request.Error(apierror.NewBadRequest("invalid auth data"))
				return
			}
			// we get the user and make sure it (still) exists
			request.User, err = auth.GetUser(sess.UserID)
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
		e.Handler(request)
	}

	return http.HandlerFunc(HTTPHandler)
}
