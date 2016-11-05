package articles_test

import "github.com/Nivl/api.melvin.la/api/app"

func init() {
	app.InitContext()
	// defer app.GetContext().Destroy()
}
