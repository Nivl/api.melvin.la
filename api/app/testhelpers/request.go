package testhelpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/api.melvin.la/api/components/api"
	"github.com/Nivl/api.melvin.la/api/router"
)

type RequestInfo struct {
	Test     *testing.T
	Endpoint *router.Endpoint
	URI      string
	Params   interface{}
}

func NewRequest(info *RequestInfo) *httptest.ResponseRecorder {
	params := bytes.NewBufferString("")

	if info.Params != nil {
		jsonDump, err := json.Marshal(info.Params)
		if err != nil {
			info.Test.Fatalf("could not create request %s", err)
		}

		params = bytes.NewBuffer(jsonDump)
	}

	req, err := http.NewRequest(info.Endpoint.Verb, info.URI, params)
	if err != nil {
		info.Test.Fatalf("could not execute request %s", err)
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	rec := httptest.NewRecorder()
	r := api.GetRouter()
	r.ServeHTTP(rec, req)
	return rec
}
