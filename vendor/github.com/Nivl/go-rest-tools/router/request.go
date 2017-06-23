package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"

	"github.com/Nivl/go-rest-tools/logger"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/gorilla/mux"
)

const (
	// ContentTypeJSON represents the content type of a JSON request
	ContentTypeJSON = "application/json"
	// ContentTypeMultipartForm represents the content type of a multipart request
	ContentTypeMultipartForm = "multipart/form-data"
)

// Request represent a client request
type Request struct {
	ID           string
	Response     http.ResponseWriter
	Request      *http.Request
	Params       interface{}
	User         *auth.User
	_contentType string
	Logger       logger.Logger
	deps         *Dependencies
}

// String return a printable version of the object
func (req *Request) String() string {
	if req == nil {
		return ""
	}

	user := "anonymous"
	userID := "0"
	if req.User != nil {
		user = req.User.Name
		userID = req.User.ID
	}

	return fmt.Sprintf(`req_id: "%s", user: "%s", user_id: "%s", endpoint: "%s", params: %#v`, req.ID, user, userID, req.Endpoint(), req.Params)
}

// ContentType returns the content type of the current request
func (req *Request) ContentType() string {
	if req == nil {
		return ""
	}

	if req._contentType == "" {
		contentType := req.Request.Header.Get("Content-Type")

		if contentType == "" {
			req._contentType = "text/html"
		} else {
			req._contentType = strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
		}
	}

	return req._contentType
}

// MuxVariables returns the URL variables associated to the request
func (req *Request) MuxVariables() url.Values {
	output := url.Values{}

	if req == nil {
		return output
	}

	vars := mux.Vars(req.Request)
	for k, v := range vars {
		output.Set(k, v)
	}

	return output
}

// JSONBody parses and returns the body of the request
func (req *Request) JSONBody() (url.Values, error) {
	output := url.Values{}

	if req.ContentType() != ContentTypeJSON {
		return output, nil
	}

	vars := map[string]interface{}{}
	if err := json.NewDecoder(req.Request.Body).Decode(&vars); err != nil && err != io.EOF {
		return nil, err
	}

	for k, v := range vars {
		output.Set(k, fmt.Sprintf("%v", v))
	}

	return output, nil
}

// ParamsBySource returns a map of params ordered by their source (url, query, form, ...)
func (req *Request) ParamsBySource() (map[string]url.Values, error) {
	params := map[string]url.Values{
		"url":   req.MuxVariables(),
		"query": req.Request.URL.Query(),
		"form":  url.Values{},
	}

	form, err := req.JSONBody()
	if err != nil {
		return nil, err
	}
	params["form"] = form

	return params, nil
}

// Endpoint returns the verb and the URI of the request
func (req *Request) Endpoint() string {
	return fmt.Sprintf("%s %s", req.Request.Method, req.Request.RequestURI)
}

// handlePanic will recover a panic an log what happen
func (req *Request) handlePanic() {
	if rec := recover(); rec != nil {
		req.Response.WriteHeader(http.StatusInternalServerError)
		req.Response.Write([]byte(`{"error":"Something went wrong"}`))

		// The recovered panic may not be an error
		var err error
		switch val := rec.(type) {
		case error:
			err = val
		default:
			err = fmt.Errorf("%v", val)
		}
		err = fmt.Errorf("panic: %v", err)

		req.Logger.Errorf(`message: "%s", %s`, err.Error(), req)
		req.Logger.Errorf(string(debug.Stack()))

		// Send an email async
		sendEmail := func(stacktrace []byte) {
			err := req.deps.Mailer.SendStackTrace(stacktrace, req.Endpoint(), err.Error(), req.ID)
			if err != nil {
				req.Logger.Error(err.Error())
			}
		}

		go sendEmail(debug.Stack())
	}
}
