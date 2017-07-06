package httptests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"strings"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/params"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/gorilla/mux"
)

// DefaultRouter represents the router to be used if none are specified
var DefaultRouter *mux.Router

// RequestAuth represents the auth data for a request
type RequestAuth struct {
	SessionID string
	UserID    string
}

// ToBasicAuth returns the data using the basic auth format
func (ra *RequestAuth) ToBasicAuth() string {
	authValue := fmt.Sprintf("%s:%s", ra.UserID, ra.SessionID)

	encoded := base64.StdEncoding.EncodeToString([]byte(authValue))
	return "basic " + encoded
}

// NewRequestAuth creates a new request auth
func NewRequestAuth(s *auth.Session) *RequestAuth {
	return &RequestAuth{
		SessionID: s.ID,
		UserID:    s.UserID,
	}
}

// RequestInfo represents the params accepted by NewRequest
type RequestInfo struct {
	Endpoint *router.Endpoint

	Params interface{}  // Optional
	Auth   *RequestAuth // Optional
	// Router is used to parse Mux Variables. Default on the api router
	Router *mux.Router

	urlParams   map[string]string
	bodyParams  map[string]string
	queryParams map[string]string
}

// ParseParams parses the params and copy them in the right list:
// urlParams, bodyParams, and queryParams
func (ri *RequestInfo) ParseParams() {
	ri.urlParams = map[string]string{}
	ri.bodyParams = map[string]string{}
	ri.queryParams = map[string]string{}

	if ri.Params == nil {
		return
	}

	p := params.NewParams(ri.Params)
	sources := p.Extract()

	for k, v := range sources["url"] {
		ri.urlParams[k] = v[0]
	}
	for k, v := range sources["form"] {
		ri.bodyParams[k] = v[0]
	}
	for k, v := range sources["query"] {
		ri.queryParams[k] = v[0]
	}
}

// URL returns the full URL
func (ri *RequestInfo) URL() string {
	url := ri.Endpoint.Path
	for param, value := range ri.urlParams {
		url = strings.Replace(url, "{"+param+"}", value, -1)
	}

	return url
}

// PopulateQuery populate the query string of a request
func (ri *RequestInfo) PopulateQuery(qs *url.Values) {
	for key, value := range ri.queryParams {
		qs.Add(key, value)
	}
}

// Body returns the full Body of the request
func (ri *RequestInfo) Body() (*bytes.Buffer, error) {
	body := bytes.NewBufferString("")

	// Parse the body as a JSON object
	if len(ri.bodyParams) > 0 {
		jsonDump, err := json.Marshal(ri.bodyParams)
		if err != nil {
			return nil, err
		}

		body = bytes.NewBuffer(jsonDump)
	}

	return body, nil
}

// NewRequest simulates a new http request executed against the api
func NewRequest(t *testing.T, info *RequestInfo) *httptest.ResponseRecorder {
	info.ParseParams()

	// Parse the body as a JSON object
	body, err := info.Body()
	if err != nil {
		t.Fatalf("could not create request: %s", err)
	}

	req, err := http.NewRequest(info.Endpoint.Verb, info.URL(), body)
	if err != nil {
		t.Fatalf("could not execute request %s", err)
	}

	// Attach the query string
	qs := req.URL.Query()
	info.PopulateQuery(&qs)
	req.URL.RawQuery = qs.Encode()

	if info.Auth != nil {
		req.Header.Add("Authorization", info.Auth.ToBasicAuth())
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// If no router is provided we assume that we want to execute a regular endpoint
	if info.Router == nil {
		if DefaultRouter == nil {
			t.Fatalf("no router provided and DefaultRouter not set")
		}
		info.Router = DefaultRouter
	}

	rec := httptest.NewRecorder()
	info.Router.ServeHTTP(rec, req)
	return rec
}
