package router

import (
	"errors"
	"reflect"
)

type RouteAuth func(*Request) bool
type RouteHandler func(*Request)

// Endpoint represents an HTTP endpoint
type Endpoint struct {
	Verb    string
	Path    string
	Auth    RouteAuth
	Handler RouteHandler
	// JSONBodyTemplate is an instance of the struct type that'll
	// be parsed out of the request body JSON.
	JSONBodyTemplate interface{}
}

// jsonBodyTemplateType returns the reflect.Type needed to create new
// instances of the JSONBodyTemplate struct.
func (e *Endpoint) jsonBodyTemplateType() (reflect.Type, error) {
	// Reflect upon the type of e.JSONBodyTemplate and make sure
	// that it's a pointer.
	payloadPtrValue := reflect.ValueOf(e.JSONBodyTemplate)
	payloadPtrType := payloadPtrValue.Type()
	if t := payloadPtrType.Kind(); t != reflect.Ptr {
		return nil, errors.New("payload definition is not a pointer")
	}

	// Follow the pointer to extract the type of the payload.
	return reflect.Indirect(payloadPtrValue).Type(), nil
}
