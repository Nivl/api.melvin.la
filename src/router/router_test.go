package router_test

import "github.com/melvin-laplanche/ml-api/src/app"

func init() {
	app.InitContext()
	// defer app.GetContext().Destroy()
}
