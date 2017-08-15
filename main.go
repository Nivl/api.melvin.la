package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/melvin-laplanche/ml-api/src/components/api"
)

func main() {
	args, err := api.Setup()
	if err != nil {
		panic(err)
	}

	r := api.GetRouter()
	port := ":" + args.Port

	handler := handlers.CORS(api.AllowedOrigins, api.AllowedMethods, api.AllowedHeaders)
	http.ListenAndServe(port, handler(r))
}
