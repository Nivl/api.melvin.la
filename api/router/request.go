package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Nivl/api.melvin.la/api/logger"
)

const CTFormData = "application/x-www-form-urlencoded"
const CTMultipartFormData = "multipart/form-data"

type Request struct {
	ID           string
	Response     http.ResponseWriter
	Request      *http.Request
	Params       url.Values
	_contentType string
}

func (req *Request) String() string {
	if req == nil {
		return ""
	}

	dump, err := json.Marshal(req)
	if err != nil {
		logger.Errorf(err.Error())
		return ""
	}

	return string(dump)
}

func (req *Request) GetContentType() string {
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
