package db

import (
	"github.com/Nivl/sqalx"
	"github.com/melvin-laplanche/ml-api/src/app"
)

// Con returns a connection to the database
func Con() sqalx.Node {
	return app.GetContext().SQL
}
