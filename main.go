package main

import (
	"net/http"

	"github.com/Nivl/cors"
	"github.com/melvin-laplanche/ml-api/src/app"
	"github.com/melvin-laplanche/ml-api/src/components/api"
)

func main() {
	app.InitContext()
	defer app.GetContext().Destroy()

	r := api.GetRouter()
	port := ":" + app.GetContext().Params.Port

	c := cors.New(cors.Options{
		AllowedOrigins: api.AllowedOrigins,
		AllowedMethods: api.AllowedMethods,
		AllowedHeaders: api.AllowedHeaders,
	})

	handler := c.Handler(r)
	http.ListenAndServe(port, handler)
}
