package router

import (
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
			Handler(Handler(endpoint))
	}
}

// Handler makes it possible to use a RouteHandler where a http.Handler is required
func Handler(endpoint *Endpoint) http.Handler {
	HTTPHandler := func(resWriter http.ResponseWriter, req *http.Request) {
		request := &Request{
			ID:       uuid.NewV4().String()[:8],
			Request:  req,
			Response: resWriter,
		}

		request.Response.Header().Set("X-Request-Id", request.ID)

		if err := endpoint.ParseParams(request); err != nil {
			request.Error(err)
			return
		}

		defer request.handlePanic()

		accessGranted := endpoint.Auth == nil || endpoint.Auth(request)
		if accessGranted {
			endpoint.Handler(request)
		}
	}

	return http.HandlerFunc(HTTPHandler)
}
