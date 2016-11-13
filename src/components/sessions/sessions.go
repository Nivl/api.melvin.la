package sessions

import (
	"github.com/melvin-laplanche/ml-api/src/app"
	"github.com/jmoiron/sqlx"
)

func sql() *sqlx.DB {
	return app.GetContext().SQL
}
