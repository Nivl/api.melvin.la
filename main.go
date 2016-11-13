package main

import (
	"net/http"

	"github.com/melvin-laplanche/ml-api/src/app"
	"github.com/melvin-laplanche/ml-api/src/components/api"
)

func main() {
	app.InitContext()
	defer app.GetContext().Destroy()

	r := api.GetRouter()
	port := ":" + app.GetContext().Params.Port
	http.ListenAndServe(port, r)
}
