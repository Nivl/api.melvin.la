package sessions

import (
	"github.com/Nivl/api.melvin.la/api/app"
	"github.com/jmoiron/sqlx"
)

func sql() *sqlx.DB {
	return app.GetContext().SQL
}
