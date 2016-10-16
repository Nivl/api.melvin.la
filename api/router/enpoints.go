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
func Handler(e *Endpoint) http.Handler {
	HTTPHandler := func(resWriter http.ResponseWriter, req *http.Request) {
		request := &Request{
			ID:       uuid.NewV4().String()[:8],
			Request:  req,
			Response: resWriter,
		}

		request.Response.Header().Set("X-Request-Id", request.ID)

		// TODO(melvin): It makes more sense to parse the params on the requests
		// we need to use reflection to create the right instance of request.Params
		// using the type of e.Params (which should only be holding a type)
		// And then doing `request.ParseParseParams()`
		// http://stackoverflow.com/questions/7850140/how-do-you-create-a-new-instance-of-a-struct-from-its-type-at-runtime-in-go
		if err := e.ParseParams(request); err != nil {
			request.Error(err)
			return
		}
		request.Params = e.Params

		defer request.handlePanic()

		accessGranted := e.Auth == nil || e.Auth(request)
		if accessGranted {
			e.Handler(request)
		}
	}

	return http.HandlerFunc(HTTPHandler)
}
