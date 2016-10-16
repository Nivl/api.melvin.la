package router

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/Nivl/api.melvin.la/api/apierror"
)

type RouteAuth func(*Request) bool
type RouteHandler func(*Request)

// Endpoint represents an HTTP endpoint
type Endpoint struct {
	Verb    string
	Path    string
	Auth    RouteAuth
	Handler RouteHandler
	Params  interface{}
}

// ParseParams will parse the params from the given request, and store them
// into the endpoint
func (e *Endpoint) ParseParams(r *Request) error {
	params := reflect.ValueOf(e.Params)
	if params.Kind() == reflect.Ptr {
		params = params.Elem()
	}

	sources, err := r.ParamsBySource()
	if err != nil {
		return err
	}

	nbParams := params.NumField()
	for i := 0; i < nbParams; i++ {
		param := params.Field(i)
		paramInfo := params.Type().Field(i)
		tags := paramInfo.Tag

		if param.Kind() == reflect.Ptr {
			param = param.Elem()
		}

		// We make sure we can update the value of field
		if !param.CanSet() {
			return apierror.NewServerError("Field %s could not be set", paramInfo.Name)
		}

		// We control the type of
		paramLocation := strings.ToLower(tags.Get("from"))
		source, found := sources[paramLocation]
		if !found {
			source = sources["url"]
		}

		args := &setParamValueArgs{
			param:     &param,
			paramInfo: &paramInfo,
			tags:      &tags,
			source:    &source,
		}

		if err := e.setParamValue(args); err != nil {
			return err
		}
	}

	return nil
}

type setParamValueArgs struct {
	param     *reflect.Value
	paramInfo *reflect.StructField
	tags      *reflect.StructTag
	source    *url.Values
}

func (e *Endpoint) setParamValue(args *setParamValueArgs) error {
	// We parse the tag to get the options
	opts := NewParamOptions(args.tags)
	defaultValue := args.tags.Get("default")

	// The tag needs to be ignored
	if opts.Ignore {
		return nil
	}

	if opts.Name == "" {
		opts.Name = args.paramInfo.Name
	}

	// We get the value and apply the transformations
	value := args.source.Get(opts.Name)
	if opts.Trim {
		value = strings.TrimSpace(value)
	}

	if value == "" {
		value = defaultValue
		if opts.Required {
			return apierror.NewBadRequest("parameter [%s] missing", opts.Name)
		}
	}

	// We now set the value in the struct
	if value != "" {
		var errorMsg = fmt.Sprintf("value [%s] for parameter [%s] is invalid", value, opts.Name)

		switch args.param.Kind() {
		case reflect.Bool:
			v, err := strconv.ParseBool(value)
			if err != nil {
				return apierror.NewBadRequest(errorMsg)
			}
			args.param.SetBool(v)
		case reflect.String:
			args.param.SetString(value)
		case reflect.Int:
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return apierror.NewBadRequest(errorMsg)
			}
			args.param.SetInt(v)
		}
	}
	return nil
}
