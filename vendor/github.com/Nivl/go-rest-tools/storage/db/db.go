package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Writer represents an open connection to the database
var Writer Transactionable

// Queryable represent a global interface for transactions, prepared statement,
// etc.
type Queryable interface {
	sqlx.Execer
	sqlx.Queryer
	sqlx.Preparer
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	DriverName() string
	Get(dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	Preparex(query string) (*sqlx.Stmt, error)
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) error
}

// Transactionable represent a global interface to create transaction
type Transactionable interface {
	Queryable
	Beginx() (*sqlx.Tx, error)
}
