package router

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

type Endpoints []*Endpoint

func (endpoints Endpoints) Activate(router *mux.Router) {
	for _, endpoint := range endpoints {
		routeHandler := endpoint.Handler

		if endpoint.JSONBodyTemplate != nil {
			bodyType, err := endpoint.getJSONBodyTemplateType()
			if err != nil {
				panic(err)
			}

			routeHandler = jsonBodyHandler(endpoint.Handler, bodyType)
		}

		router.
			Methods(endpoint.Verb).
			Path(endpoint.Path).
			Handler(Handler(routeHandler, endpoint.Auth))
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

func jsonBodyHandler(next RouteHandler, pldType reflect.Type) RouteHandler {
	return func(req *Request) {
		pld := reflect.New(pldType).Interface()

		if req.Request.Body == nil {
			req.BadRequest("JSON body not provided")
			return
		}

		if err := json.NewDecoder(req.Request.Body).Decode(pld); err != nil {
			req.BadRequest("could not parse JSON body")
			return
		}

		req.JSONBody = pld

		next(req)
	}
}
