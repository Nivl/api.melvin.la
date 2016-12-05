package testhelpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/melvin-laplanche/ml-api/src/components/api"
	"github.com/melvin-laplanche/ml-api/src/router"
)

type RequestAuth struct {
	SessionID string
	UserID    string
}

func NewRequestAuth(sessionID string, userID string) *RequestAuth {
	return &RequestAuth{
		SessionID: sessionID,
		UserID:    userID,
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
		req.Header.Add("X-Session-Token", info.Auth.SessionID)
		req.Header.Add("X-User-Id", info.Auth.UserID)
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

// Is2XX returns a true if the HTTP code is a 2XX, false otherwise
func Is2XX(code int) bool {
	return code >= 200 && code < 300
}
