package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/Nivl/api.melvin.la/api/logger"
)

func (req *Request) Error(e error) {
	if req == nil {
		return
	}

	// todo(melvin): can we just cast e into `apierror.Error` and check
	// if err.Code == 0
	var err apierror.Error
	switch e.(type) {
	case apierror.Error:
		err = e.(apierror.Error)
	default:
		err = apierror.NewServerError(err.Error()).(apierror.Error)
	}

	switch err.Code {
	case http.StatusInternalServerError:
		logger.Errorf("%s - %s", err.Error(), req)
		http.Error(req.Response, `{"error":"Something went wrong"}`, http.StatusInternalServerError)
	default:
		http.Error(req.Response, fmt.Sprintf(`{"error":"%s"}`, err.Error()), err.Code)
	}
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
			req.Error(err)
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
