package router

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/serenize/snaker"
)

type RouteAuth func(*Request) bool
type RouteHandler func(*Request)

// Endpoint represents an HTTP endpoint
type Endpoint struct {
	Verb    string
	Path    string
	Auth    RouteAuth
	Handler RouteHandler
	Params  Params
}

func (e *Endpoint) ParseParams(r *Request) error {
	// Query String (URL)
	if err := e.extractParams(&e.Params.Query, r.Request.URL.Query()); err != nil {
		return err
	}

	// Form Data
	if r.GetContentType() == CTFormData {
		if err := r.Request.ParseForm(); err != nil {
			return apierror.NewServerError("Error parsing request body %s", err.Error())
		}

		if err := e.extractParams(&e.Params.Form, r.Request.Form); err != nil {
			return err
		}
	}
	return nil
}

// extractParams returns a subset of `params` that matched `definitions`
func (e *Endpoint) extractParams(definitions *interface{}, params url.Values) error {
	wantedParams := reflect.ValueOf(definitions).Elem()
	nbWantedParams := wantedParams.NumField()

	for i := 0; i < nbWantedParams; i++ {
		wantedParam := wantedParams.Field(i)
		wantedParamInfo := wantedParams.Type().Field(i)
		tags := wantedParamInfo.Tag

		// We make sure we can update the value of field
		if !wantedParam.CanSet() {
			return apierror.NewServerError("Field %s could not be set", wantedParamInfo.Name)
		}

		// We parse the tag to get the options
		opts := NewParamOptions(&tags)
		defaultValue := tags.Get("default")

		// If no name has been specified, we'll use the snake version
		// of the variable name
		if opts.Name == "" {
			opts.Name = snaker.SnakeToCamel(wantedParamInfo.Name)
		}

		// We get the valye apply the transformations
		value := params.Get(opts.Name)
		if opts.Trim {
			value = strings.TrimSpace(value)
		}

		if value == "" {
			value = defaultValue
			if opts.Required {
				return apierror.NewBadRequest("parameter '%s' missing", opts.Name)
			}
		}

		// We now set the value in the struct
		if value != "" {
			switch wantedParam.Kind() {
			case reflect.Bool:
				v, err := strconv.ParseBool(value)
				if err != nil {
					return apierror.NewBadRequest("parameter '%s' missing", opts.Name)
				}
				wantedParam.SetBool(v)
			case reflect.String:
				wantedParam.SetString(value)
			case reflect.Int:
				v, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return apierror.NewBadRequest("parameter '%s' missing", opts.Name)
				}
				wantedParam.SetInt(v)
			}
		}
	}

	return nil
}
