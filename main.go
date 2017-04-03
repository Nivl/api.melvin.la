package main

import (
	"net/http"

	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/gorilla/handlers"
	"github.com/melvin-laplanche/ml-api/src/components/api"
)

func main() {
	args := api.Setup()
	defer db.Writer.Close()

	r := api.GetRouter()
	port := ":" + args.Port

	handler := handlers.CORS(api.AllowedOrigins, api.AllowedMethods, api.AllowedHeaders)
	http.ListenAndServe(port, handler(r))
}
