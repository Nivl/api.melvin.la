package app

import (
	"fmt"
	"github.com/bsphere/le_go"
	"github.com/jessevdk/go-flags"
	"gopkg.in/mgo.v2"
)

// Args represents the app args
type Args struct {
	Port            string `short:"p" long:"port" description:"Port to listen on" default:"8000"`
	MongoURI        string `long:"database" description:"Connection URI of the mongo database" required:"true"`
	LogEntriesToken string `long:"logentries" description:"Logenentries token to log errors"`
	Debug           bool   `short:"d" long:"debug" description:"Debug mode"`
}

// Context represent the global context of the app
type Context struct {
	DB         *mgo.Database
	Session    *mgo.Session
	Params     Args
	LogEntries *le_go.Logger
}

var _context *Context

// InitContextWithParams initializes the app context using the provided params
func InitContextWithParams(argv *Args) *Context {
	if _context != nil {
		panic("Context already exists")
	}

	_context = new(Context)

	// Parse the cmd args
	if argv == nil {
		_, err := flags.Parse(&_context.Params)
		if err != nil {
			panic(err)
		}
	} else {
		_context.Params = *argv
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

// InitContex initializes the app context using argv
func InitContex() *Context {
	return InitContextWithParams(nil)
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
