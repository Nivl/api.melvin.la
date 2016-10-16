package router

import (
	"reflect"
	"strings"
)

/*
	Accept the folowing tags
	require


	Example:

	{
		Title        string `params:",required"`
		IsPublished  string `params:"" default:""`
		Private  string `params:"-"`
	}
*/

type ParamOptions struct {
	Ignore   bool
	Name     string
	Required bool
	Trim     bool
}

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
		}
	}

	return output
}
