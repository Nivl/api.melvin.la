package testhelpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/api.melvin.la/api/components/api"
	"github.com/Nivl/api.melvin.la/api/router"
	"github.com/gorilla/mux"
)

type RequestAuth struct {
	SessionUUID string
	UserUUID    string
}

func NewRequestAuth(sessionUUID string, userUUID string) *RequestAuth {
	return &RequestAuth{
		SessionUUID: sessionUUID,
		UserUUID:    userUUID,
	}
}

type RequestInfo struct {
	Test     *testing.T
	Endpoint *router.Endpoint
	URI      string
	Params   interface{}
	Auth     *RequestAuth

	// Router is used to parse Mux Variables. Default on the api router
	Router *mux.Router
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

	if info.Auth != nil {
		req.Header.Add("X-Session-Token", info.Auth.SessionUUID)
		req.Header.Add("X-User-Id", info.Auth.UserUUID)
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// If no router is provided we assume that we want to execute a regular endpoint
	if info.Router == nil {
		info.Router = api.GetRouter()
	}

	rec := httptest.NewRecorder()
	info.Router.ServeHTTP(rec, req)
	return rec
}
