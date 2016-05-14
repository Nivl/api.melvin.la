package app

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"gopkg.in/mgo.v2"
)

// Args represents the app args
type Args struct {
	Port     string `short:"p" long:"port" description:"Port to listen on"`
	MongoURI string `long:"database" description:"Connection URI of the mongo database"`
}

// Context represent the global context of the app
type Context struct {
	DB      *mgo.Database
	Session *mgo.Session
	Params  Args
}

var _context *Context

// GetContext return the current app context.
func GetContext() *Context {
	if _context == nil {
		_context = new(Context)

		_, err := flags.Parse(&_context.Params)
		if err != nil {
			panic(err)
		}

		session, err := mgo.Dial(_context.Params.MongoURI)
		if err != nil {
			fmt.Println("Cannot start mongo")
			panic(err)
		}

		_context.Session = session
		_context.Session.SetMode(mgo.Monotonic, true)
		_context.DB = session.DB("")
	}

	return _context
}

// Destroy clears what's need to be cleared when the app is quiting
func (ctx *Context) Destroy() {
	if ctx.Session != nil {
		ctx.Session.Close()
	}
}
