package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type Endpoints []*Endpoint

func (endpoints Endpoints) Activate(router *mux.Router) {
	for _, endpoint := range endpoints {
		router.
			Methods(endpoint.Verb).
			Path(endpoint.Path).
			Handler(Handler(endpoint.Handler, endpoint.Auth))
	}
}

// Handler makes it possible to use a RouteHandler where a http.Handler is required
func Handler(handler RouteHandler, auth RouteAuth) http.Handler {
	HTTPHandler := func(resWriter http.ResponseWriter, req *http.Request) {
		request := &Request{
			ID:       uuid.NewV4().String()[:8],
			Request:  req,
			Response: resWriter,
		}

		request.Response.Header().Set("X-Request-Id", request.ID)
		// TODO: handle users

		err := request.ParseParams()
		defer removeParams(request)
		// We must return a 400 and stop here if there was a problem parsing the request.
		if err != nil {
			http.Error(request.Response, fmt.Sprintf(`{"error":"%s"}`, "Bad params"), 400)
			return
		}

		defer request.handlePanic()

		accessGranted := auth == nil || auth(request)
		if accessGranted {
			handler(request)
		}
	}

	return http.HandlerFunc(HTTPHandler)
}
