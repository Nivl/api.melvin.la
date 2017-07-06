package params

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
)

type Params struct {
	data interface{}
}

func NewParams(data interface{}) *Params {
	return &Params{
		data: data,
	}
}

// Parse fills the paramsStruct using the provided sources
func (p *Params) Parse(sources map[string]url.Values) error {
	paramList := reflect.ValueOf(p.data)
	if paramList.Kind() == reflect.Ptr {
		paramList = paramList.Elem()
	}

	return p.parseRecursive(paramList, sources)
}

func (p *Params) parseRecursive(paramList reflect.Value, sources map[string]url.Values) error {
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
			p.parseRecursive(value, sources)
			continue
		}

		// We control the source of the param. If nothing is provided, we take from the URL
		paramLocation := strings.ToLower(tags.Get("from"))
		if paramLocation == "" {
			paramLocation = "url"
		}

		source, found := sources[paramLocation]
		if !found {
			return httperr.NewServerError("source [%s] for field [%s] does not exists", paramLocation, info.Name)
		}

		param := &Param{
			value: &value,
			info:  &info,
			tags:  &tags,
		}
		if err := param.SetValue(&source); err != nil {
			return err
		}
	}

	return nil
}

// Extract extracts the data from the paramsStruct and returns them
// as a map of url.Values
func (p *Params) Extract() map[string]url.Values {
	sources := map[string]url.Values{}
	sources["url"] = url.Values{}
	sources["form"] = url.Values{}
	sources["query"] = url.Values{}

	if p.data == nil {
		return sources
	}

	paramList := reflect.ValueOf(p.data)
	if paramList.Kind() == reflect.Ptr {
		paramList = paramList.Elem()
	}
	p.extractRecursive(paramList, sources)
	return sources
}

func (p *Params) extractRecursive(paramList reflect.Value, sources map[string]url.Values) {
	nbParams := paramList.NumField()
	for i := 0; i < nbParams; i++ {
		value := paramList.Field(i)
		paramInfo := paramList.Type().Field(i)
		tags := paramInfo.Tag

		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}

		// Handle embedded struct
		if value.Kind() == reflect.Struct && paramInfo.Anonymous {
			p.extractRecursive(value, sources)
			continue
		}

		// We get the Value as string
		valueStr := ""
		switch value.Kind() {
		case reflect.Bool:
			valueStr = strconv.FormatBool(value.Bool())
		case reflect.String:
			valueStr = value.String()
		case reflect.Int:
			valueStr = strconv.Itoa(int(value.Int()))
		case reflect.Ptr:
			if !value.IsNil() {
				val := value.Elem()
				switch val.Kind() {
				case reflect.Bool:
					valueStr = strconv.FormatBool(val.Bool())
				case reflect.String:
					valueStr = value.String()
				case reflect.Int:
					valueStr = strconv.Itoa(int(val.Int()))
				}
			}
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
		sources[sourceType].Set(fieldName, valueStr)
	}
}
