package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/melvin-laplanche/ml-api/src/components/api"
)

func main() {
	params, deps, err := api.DefaultSetup()
	if err != nil {
		panic(err)
	}

	r := api.GetRouter(deps)
	port := ":" + params.Port

	handler := handlers.CORS(api.AllowedOrigins, api.AllowedMethods, api.AllowedHeaders)
	http.ListenAndServe(port, handler(r))
}
