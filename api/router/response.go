package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nivl/api.melvin.la/api/logger"
)

func (req *Request) ServerError(err error) {
	if req == nil {
		return
	}

	logger.Errorf("%s - %s", err.Error(), req)
	http.Error(req.Response, `{"error":"Something went wrong"}`, http.StatusInternalServerError)
}

func (req *Request) NotFound(msg string, args ...interface{}) {
	if req == nil {
		return
	}

	fullMsg := fmt.Sprintf(msg, args...)
	http.Error(req.Response, fmt.Sprintf(`{"error":"%s"}`, fullMsg), http.StatusNotFound)
}

func (req *Request) Conflict(msg string, args ...interface{}) {
	if req == nil {
		return
	}

	fullMsg := fmt.Sprintf(msg, args...)
	http.Error(req.Response, fmt.Sprintf(`{"error":"%s"}`, fullMsg), http.StatusConflict)
}

func (req *Request) BadRequest(msg string, args ...interface{}) {
	if req == nil {
		return
	}

	fullMsg := fmt.Sprintf(msg, args...)
	http.Error(req.Response, fmt.Sprintf(`{"error":"%s"}`, fullMsg), http.StatusBadRequest)
}

func (req *Request) NoContent() {
	if req == nil {
		return
	}

	req.Response.WriteHeader(http.StatusNoContent)
}

func (req *Request) Created(obj interface{}) {
	if req == nil {
		return
	}

	req.RenderJSON(http.StatusCreated, obj)
}

func (req *Request) Ok(obj interface{}) {
	if req == nil {
		return
	}

	req.RenderJSON(http.StatusOK, obj)
}

func (req *Request) RenderJSON(code int, obj interface{}) {
	var err error
	var dump []byte

	if obj != nil {
		dump, err = json.Marshal(obj)
		if err != nil {
			req.ServerError(err)
			return
		}
	}

	req.Response.WriteHeader(code)

	if len(dump) == 0 {
		if _, err = req.Response.Write(dump); err != nil {
			req.Response.WriteHeader(http.StatusInternalServerError)
			logger.Errorf("Could not write JSON response: %s", err.Error())
		}
	}
}
