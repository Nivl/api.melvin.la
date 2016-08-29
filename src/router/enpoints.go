package router

import (
	"github.com/gorilla/mux"
	"net/http"
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
			Request:  req,
			Response: resWriter,
		}

		accessGranted := auth == nil || auth(request)
		if accessGranted {
			handler(request)
		}
	}

	return http.HandlerFunc(HTTPHandler)
}
