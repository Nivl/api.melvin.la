package test

import (
	"net/http"

	"github.com/Nivl/api.melvin.la/src/app"
)

const _endpoint = "http://localhost"

func getDomain() string {
	return _endpoint + ":" + app.GetContext().Params.Port
}

// Get returns a GET response
func Get(path string) *http.Response {
	endpoint := getDomain() + path

	response, err := http.Get(endpoint)
	if err != nil {
		panic(err)
	}

	return response
}
