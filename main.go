package main

import (
	"net/http"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/gorilla/handlers"
	"github.com/kelseyhightower/envconfig"
	"github.com/melvin-laplanche/ml-api/src/components/api"
)

func main() {
	params := &api.Args{}
	if err := envconfig.Process("", params); err != nil {
		panic(err)
	}

	deps := &dependencies.APIDependencies{}
	err := api.Setup(params, deps)
	if err != nil {
		panic(err)
	}

	r := api.GetRouter(deps)
	port := ":" + params.Port

	handler := handlers.CORS(api.AllowedOrigins, api.AllowedMethods, api.AllowedHeaders)
	http.ListenAndServe(port, handler(r))
}
