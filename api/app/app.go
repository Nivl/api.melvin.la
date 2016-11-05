package app

import (
	"fmt"

	"github.com/bsphere/le_go"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/mgo.v2"
)

// Args represents the app args
type Args struct {
	Port            string `default:"5000"`
	MongoURI        string `required:"true" envconfig:"mongo_uri"`
	LogEntriesToken string `envconfig:"mongo_uri" envconfig:"logentries_token"`
	Debug           bool   `default:"false"`
}

// Context represent the global context of the app
type Context struct {
	DB         *mgo.Database
	Session    *mgo.Session
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

	if err := envconfig.Process("api", &_context.Params); err != nil {
		panic(err)
	}

	// Setup database
	session, err := mgo.Dial(_context.Params.MongoURI)
	if err != nil {
		fmt.Println("Cannot start mongo")
		panic(err)
	}
	_context.Session = session
	_context.Session.SetMode(mgo.Monotonic, true)
	_context.DB = session.DB("")

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

// Destroy clears the context when the app is quiting
func (ctx *Context) Destroy() {
	if ctx.Session != nil {
		ctx.Session.Close()
	}

	if ctx.LogEntries != nil {
		ctx.LogEntries.Close()
	}
}
