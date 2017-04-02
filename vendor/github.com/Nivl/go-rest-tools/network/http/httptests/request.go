package httptests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"strings"

	"github.com/Nivl/go-rest-tools/router"
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

// NewRequestAuth creates a new
func NewRequestAuth(sessionID string, userID string) *RequestAuth {
	return &RequestAuth{
		SessionID: sessionID,
		UserID:    userID,
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

	params := reflect.ValueOf(ri.Params)
	if params.Kind() == reflect.Ptr {
		params = params.Elem()
	}
	ri.parseParamsRecursive(params)
}

func (ri *RequestInfo) parseParamsRecursive(params reflect.Value) {
	nbParams := params.NumField()
	for i := 0; i < nbParams; i++ {
		param := params.Field(i)
		paramInfo := params.Type().Field(i)
		tags := paramInfo.Tag

		if param.Kind() == reflect.Ptr {
			param = param.Elem()
		}

		// Handle embeded struct
		if param.Kind() == reflect.Struct && paramInfo.Anonymous {
			ri.parseParamsRecursive(param)
			continue
		}

		// We get the Value
		value := ""
		switch param.Kind() {
		case reflect.Bool:
			value = strconv.FormatBool(param.Bool())
		case reflect.String:
			value = param.String()
		case reflect.Int:
			value = strconv.Itoa(int(param.Int()))
		case reflect.Ptr:
			if !param.IsNil() {
				val := param.Elem()
				switch val.Kind() {
				case reflect.Bool:
					value = strconv.FormatBool(val.Bool())
				case reflect.String:
					value = param.String()
				case reflect.Int:
					value = strconv.Itoa(int(val.Int()))
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

		switch strings.ToLower(tags.Get("from")) {
		case "url":
			ri.urlParams[fieldName] = value
		case "form":
			ri.bodyParams[fieldName] = value
		case "query":
			ri.queryParams[fieldName] = value
		}
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
		info.Router = DefaultRouter
	}

	rec := httptest.NewRecorder()
	info.Router.ServeHTTP(rec, req)
	return rec
}
