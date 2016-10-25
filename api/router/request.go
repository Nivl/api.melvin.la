package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/logger"
	"github.com/gorilla/mux"
)

const (
	ContentTypeJSON          = "application/json"
	ContentTypeMultipartForm = "multipart/form-data"
)

type Request struct {
	ID           string              `json:"req_id"`
	Response     http.ResponseWriter `json:"-"`
	Request      *http.Request       `json:"-"`
	Params       interface{}         `json:"-"`
	User         *auth.User          `json:"-"`
	_contentType string
}

func (req *Request) String() string {
	if req == nil {
		return ""
	}

	dump, err := json.Marshal(req)
	if err != nil {
		logger.Errorf(err.Error())
		return "failed to parse the request"
	}

	return string(dump)
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
	var output url.Values

	if req == nil {
		return output
	}

	vars := mux.Vars(req.Request)
	for k, v := range vars {
		output.Set(k, v)
	}

	return output
}

// MuxVariables parses and returns the body of the request
func (req *Request) JSONBody() (url.Values, error) {
	output := url.Values{}

	if req.ContentType() != ContentTypeJSON {
		return output, nil
	}

	vars := map[string]string{}
	if err := json.NewDecoder(req.Request.Body).Decode(&vars); err != nil {
		return nil, err
	}

	for k, v := range vars {
		output.Set(k, v)
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
		// TODO send an email
		logger.Error(err.Error())
	}
}
