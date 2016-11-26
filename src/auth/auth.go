package auth

import (
	"github.com/jmoiron/sqlx"
	"github.com/melvin-laplanche/ml-api/src/app"
)

func sql() *sqlx.DB {
	return app.GetContext().SQL
}
