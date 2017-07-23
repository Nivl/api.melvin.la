package params

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router/formfile"
)

// Params is a struct used to parse and extract params from an other struct
type Params struct {
	data interface{}
}

// NewParams creates a new Params object from a struct
func NewParams(data interface{}) *Params {
	return &Params{
		data: data,
	}
}

// Parse fills the paramsStruct using the provided sources
func (p *Params) Parse(sources map[string]url.Values, fileHolder formfile.FileHolder) error {
	paramList := reflect.Indirect(reflect.ValueOf(p.data))
	return p.parseRecursive(paramList, sources, fileHolder)
}

func (p *Params) parseRecursive(paramList reflect.Value, sources map[string]url.Values, fileHolder formfile.FileHolder) error {
	nbParams := paramList.NumField()
	for i := 0; i < nbParams; i++ {
		value := paramList.Field(i)
		info := paramList.Type().Field(i)
		tags := info.Tag

		// We make sure we can update the value of field
		if !value.CanSet() {
			return httperr.NewServerError("field [%s] could not be set", info.Name)
		}

		// Handle embedded struct
		if value.Kind() == reflect.Struct && info.Anonymous {
			p.parseRecursive(value, sources, fileHolder)
			continue
		}

		// We control the source of the param. If nothing is provided, we take from the URL
		paramLocation := strings.ToLower(tags.Get("from"))
		if paramLocation == "" {
			paramLocation = "url"
		}

		param := &Param{
			value: &value,
			info:  &info,
			tags:  &tags,
		}

		// the "file" source is a special case as it's not part of the sources object
		if paramLocation == "file" {
			if err := param.SetFile(fileHolder); err != nil {
				return err
			}
		} else {
			source, found := sources[paramLocation]
			if !found {
				return httperr.NewServerError("source [%s] for field [%s] does not exists", paramLocation, info.Name)
			}

			if err := param.SetValue(source); err != nil {
				return err
			}
		}
	}

	return nil
}

// Extract extracts the data from the paramsStruct and returns them
// as a map of url.Values
func (p *Params) Extract() (map[string]url.Values, map[string]*formfile.FormFile) {
	sources := map[string]url.Values{}
	sources["url"] = url.Values{}
	sources["form"] = url.Values{}
	sources["query"] = url.Values{}
	files := map[string]*formfile.FormFile{}

	if p.data == nil {
		return sources, files
	}

	paramList := reflect.Indirect(reflect.ValueOf(p.data))
	p.extractRecursive(paramList, sources, files)
	return sources, files
}

func (p *Params) extractRecursive(paramList reflect.Value, sources map[string]url.Values, files map[string]*formfile.FormFile) {
	nbParams := paramList.NumField()
	for i := 0; i < nbParams; i++ {
		value := paramList.Field(i)
		info := paramList.Type().Field(i)
		tags := info.Tag

		// Handle embedded struct
		if reflect.Indirect(value).Kind() == reflect.Struct && info.Anonymous {
			p.extractRecursive(value, sources, files)
			continue
		}

		// We get the name from the json tag
		fieldName := ""
		jsonOpts := strings.Split(tags.Get("json"), ",")
		if len(jsonOpts) > 0 {
			if jsonOpts[0] == "-" {
				continue
			}
			fieldName = jsonOpts[0]
		}

		// We get the source type (url, query, form, ...) and add the value
		sourceType := strings.ToLower(tags.Get("from"))

		if _, found := sources[sourceType]; !found {
			sources[sourceType] = url.Values{}
		}

		// Special cases for files
		if info.Type.String() == "*formfile.FormFile" {
			if !value.IsNil() {
				files[fieldName] = value.Interface().(*formfile.FormFile)
			}
			continue
		}

		field := reflect.Indirect(value)
		valueStr := ""
		switch field.Kind() {
		case reflect.Bool:
			valueStr = strconv.FormatBool(field.Bool())
		case reflect.String:
			valueStr = field.String()
		case reflect.Int:
			valueStr = strconv.Itoa(int(field.Int()))
		}

		sources[sourceType].Set(fieldName, valueStr)
	}
}
