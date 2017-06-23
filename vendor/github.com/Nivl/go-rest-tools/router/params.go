package router

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/primitives/strngs"
)

// ParamOptions represent all the options for a field
type ParamOptions struct {
	// Ignore means the field should not been parsed
	// json:"-"
	Ignore bool

	// Name contains the name of the field in the payload
	// json:"my_field"
	Name string

	// Required means the request should fail with a Bad Request if the field is missing.
	// params:"required"
	Required bool

	// Trim means the field needs to be trimmed before being retrieved and checked
	// params:"trim"
	Trim bool

	// ValidateUUID means the field should contain a valid UUIDv4
	// params:"uuid"
	ValidateUUID bool

	// ValidateOptionalBool means the field should either be empty or contain a bool
	// params:"bool"
	ValidateOptionalBool bool
}

// Validate checks the given value passes the options set
func (opts *ParamOptions) Validate(value string) error {
	if value == "" && opts.Required {
		return httperr.NewBadRequest("parameter missing: %s", opts.Name)
	}

	if value != "" {
		if opts.ValidateUUID && !strngs.IsValidUUID(value) {
			return httperr.NewBadRequest("not a valid uuid: %s - %s", opts.Name, value)
		}

		if opts.ValidateOptionalBool {
			if _, err := strconv.ParseBool(value); err != nil {
				return httperr.NewBadRequest("not a valid bool: %s - %s", opts.Name, value)
			}
		}
	}

	return nil
}

// ApplyTransformations applies all the wanted transformations to the given value
func (opts *ParamOptions) ApplyTransformations(value string) string {
	if opts.Trim {
		value = strings.TrimSpace(value)
	}
	return value
}

// NewParamOptions returns a ParamOptions from a StructTag
func NewParamOptions(tags *reflect.StructTag) *ParamOptions {
	output := &ParamOptions{}

	// We use the json tag to get the field name
	jsonOpts := strings.Split(tags.Get("json"), ",")
	if len(jsonOpts) > 0 {
		if jsonOpts[0] == "-" {
			return &ParamOptions{Ignore: true}
		}

		output.Name = jsonOpts[0]
	}

	// We parse the params
	opts := strings.Split(tags.Get("params"), ",")
	nbOptions := len(opts)
	for i := 0; i < nbOptions; i++ {
		switch opts[i] {
		case "required":
			output.Required = true
		case "trim":
			output.Trim = true
		case "uuid":
			output.ValidateUUID = true
		case "bool":
			output.ValidateOptionalBool = true
		}
	}

	return output
}

// ParseParams will parse the params from the given request, and store them
// into the endpoint
func (r *Request) ParseParams() error {
	params := reflect.ValueOf(r.Params)
	if params.Kind() == reflect.Ptr {
		params = params.Elem()
	}

	sources, err := r.ParamsBySource()
	if err != nil {
		return err
	}

	return r.parseParamsRecursive(params, sources)
}

func (r *Request) parseParamsRecursive(params reflect.Value, sources map[string]url.Values) error {
	nbParams := params.NumField()
	for i := 0; i < nbParams; i++ {
		param := params.Field(i)
		paramInfo := params.Type().Field(i)
		tags := paramInfo.Tag

		// We make sure we can update the value of field
		if !param.CanSet() {
			return httperr.NewServerError("field [%s] could not be set", paramInfo.Name)
		}

		// Handle embeded struct
		if param.Kind() == reflect.Struct && paramInfo.Anonymous {
			r.parseParamsRecursive(param, sources)
			continue
		}

		// We control the source of the param. If nothing is provided, we take from the URL
		paramLocation := strings.ToLower(tags.Get("from"))
		if paramLocation == "" {
			paramLocation = "url"
		}

		source, found := sources[paramLocation]
		if !found {
			return httperr.NewServerError("source [%s] for field [%s] does not exists", paramLocation, paramInfo.Name)
		}

		args := &setParamValueArgs{
			param:     &param,
			paramInfo: &paramInfo,
			tags:      &tags,
			source:    &source,
		}

		if err := r.setParamValue(args); err != nil {
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

func (r *Request) setParamValue(args *setParamValueArgs) error {
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

	value := opts.ApplyTransformations(args.source.Get(opts.Name))
	if value == "" {
		value = defaultValue
	}

	if err := opts.Validate(value); err != nil {
		return err
	}

	// We now set the value in the struct
	if value != "" {
		var errorMsg = fmt.Sprintf("value [%s] for parameter [%s] is invalid", value, opts.Name)

		switch args.param.Kind() {
		case reflect.Bool:
			v, err := strconv.ParseBool(value)
			if err != nil {
				return httperr.NewBadRequest(errorMsg)
			}
			args.param.SetBool(v)
		case reflect.String:
			args.param.SetString(value)
		case reflect.Int:
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return httperr.NewBadRequest(errorMsg)
			}
			args.param.SetInt(v)
		case reflect.Ptr:
			val := reflect.New(args.param.Type().Elem())
			args.param.Set(val)

			switch args.param.Elem().Kind() {
			case reflect.Int:
				v64, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return httperr.NewBadRequest(errorMsg)
				}
				v := int(v64)
				args.param.Set(reflect.ValueOf(&v))
			case reflect.Bool:
				v, err := strconv.ParseBool(value)
				if err != nil {
					return httperr.NewBadRequest(errorMsg)
				}
				args.param.Set(reflect.ValueOf(&v))
			case reflect.String:
				args.param.Set(reflect.ValueOf(&value))
			}
		}
	}
	return nil
}
