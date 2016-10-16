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

	opts := strings.Split(tags.Get("params"), ",")
	nbOptions := len(opts)

	if nbOptions > 0 {
		if opts[0] == "-" {
			return &ParamOptions{Ignore: true}
		}

		// We only accept the "required" option so we just do a simple check
		output.Name = opts[0]

		for i := 1; i < nbOptions; i++ {
			switch opts[i] {
			case "required":
				output.Required = true
			case "trim":
				output.Trim = true
			}
		}
	}

	return output
}
