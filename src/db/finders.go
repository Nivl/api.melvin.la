package db

import (
	"github.com/melvin-laplanche/ml-api/src/app"
	"github.com/jmoiron/sqlx"
)

func sequel() *sqlx.DB {
	return app.GetContext().SQL
}

// Get is the same as sqlx.Get() but do not returns an error on empty results
func Get(dest interface{}, query string, args ...interface{}) error {
	err := sequel().Get(dest, query, args...)

	if SQLNotFound(err) {
		return nil
	}

	return err
}
