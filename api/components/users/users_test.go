package users_test

import (
	"github.com/Nivl/api.melvin.la/api/app"
	"github.com/jmoiron/sqlx"
)

func init() {
	app.InitContext()
	// defer app.GetContext().Destroy()
}

func sql() *sqlx.DB {
	return app.GetContext().SQL
}
