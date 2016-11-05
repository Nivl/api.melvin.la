package main

import (
	"net/http"

	"github.com/Nivl/api.melvin.la/api/app"
	"github.com/Nivl/api.melvin.la/api/components/api"
)

func main() {
	app.InitContext()
	defer app.GetContext().Destroy()

	api.EnsureIndexes()
	r := api.GetRouter()
	port := ":" + app.GetContext().Params.Port
	http.ListenAndServe(port, r)
}
