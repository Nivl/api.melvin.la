package app

import (
	"github.com/Nivl/sqalx"
	"github.com/bsphere/le_go"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"

	// Required to connect to postgres
	_ "github.com/lib/pq"
)

// Args represents the app args
type Args struct {
	Port            string `default:"5000"`
	PostgresURI     string `required:"true" envconfig:"postgres_uri"`
	LogEntriesToken string `envconfig:"logentries_token"`
	Debug           bool   `default:"false"`
}

// Context represent the global context of the app
type Context struct {
	SQL        sqalx.Node
	Params     Args
	LogEntries *le_go.Logger
}

var _context *Context

// InitContext initializes the app context
func InitContext() *Context {
	if _context != nil {
		panic("Context already exists")
	}

	_context = new(Context)
	if err := envconfig.Process("", &_context.Params); err != nil {
		panic(err)
	}

	// Setup database
	db, err := sqlx.Connect("postgres", _context.Params.PostgresURI)
	if err != nil {
		panic(err)
	}

	_context.SQL, err = sqalx.New(db, sqalx.SavePoint(true))
	if err != nil {
		panic(err)
	}

	// LogEntries
	if _context.Params.LogEntriesToken != "" {
		_context.LogEntries, err = le_go.Connect(_context.Params.LogEntriesToken)
		if err != nil {
			panic(err)
		}
	}

	return _context
}

// GetContext returns the current app context.
func GetContext() *Context {
	return _context
}

// Destroy clears the context when the app is quitting
func (ctx *Context) Destroy() {
	if ctx.SQL != nil {
		ctx.SQL.Close()
	}

	if ctx.LogEntries != nil {
		ctx.LogEntries.Close()
	}
}
