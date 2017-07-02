package params

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
)

// Param represents a struct param
type Param struct {
	value  *reflect.Value
	info   *reflect.StructField
	tags   *reflect.StructTag
	source *url.Values
}

// NewParamFromStructValue creates a param using a struct value
func NewParamFromStructValue(paramList *reflect.Value, paramPos int) *Param {
	value := paramList.Field(paramPos)
	info := paramList.Type().Field(paramPos)
	tags := info.Tag

	return &Param{
		value: &value,
		info:  &info,
		tags:  &tags,
	}
}

func (p *Param) SetValue(source *url.Values) error {
	// We parse the tag to get the options
	opts := NewParamOptions(p.tags)
	defaultValue := p.tags.Get("default")

	// The tag needs to be ignored
	if opts.Ignore {
		return nil
	}

	if opts.Name == "" {
		opts.Name = p.info.Name
	}

	value := opts.ApplyTransformations(source.Get(opts.Name))
	if value == "" {
		value = defaultValue
	}

	if err := opts.Validate(value); err != nil {
		return err
	}

	// We now set the value in the struct
	if value != "" {
		var errorMsg = fmt.Sprintf("value [%s] for parameter [%s] is invalid", value, opts.Name)

		switch p.value.Kind() {
		case reflect.Bool:
			v, err := strconv.ParseBool(value)
			if err != nil {
				return httperr.NewBadRequest(errorMsg)
			}
			p.value.SetBool(v)
		case reflect.String:
			p.value.SetString(value)
		case reflect.Int:
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return httperr.NewBadRequest(errorMsg)
			}
			p.value.SetInt(v)
		case reflect.Ptr:
			val := reflect.New(p.value.Type().Elem())
			p.value.Set(val)

			switch p.value.Elem().Kind() {
			case reflect.Int:
				v64, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return httperr.NewBadRequest(errorMsg)
				}
				v := int(v64)
				p.value.Set(reflect.ValueOf(&v))
			case reflect.Bool:
				v, err := strconv.ParseBool(value)
				if err != nil {
					return httperr.NewBadRequest(errorMsg)
				}
				p.value.Set(reflect.ValueOf(&v))
			case reflect.String:
				p.value.Set(reflect.ValueOf(&value))
			}
		}
	}
	return nil
}
