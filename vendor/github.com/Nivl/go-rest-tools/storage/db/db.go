package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// DB is an interface representing a database connection
type DB interface {
	sqlx.Ext
	sqlx.Preparer
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}
