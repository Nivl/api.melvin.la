package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/app"
	"github.com/melvin-laplanche/ml-api/src/logger"
)

func (req *Request) Error(e error) {
	if req == nil {
		return
	}

	err, casted := e.(*apierror.ApiError)
	if !casted {
		err = apierror.NewServerError(e.Error()).(*apierror.ApiError)
	}

	switch err.Code() {
	case http.StatusInternalServerError:
		logger.Errorf("%s - %s", err.Error(), req)
		http.Error(req.Response, `{"error":"Something went wrong"}`, http.StatusInternalServerError)
	default:
		// Some errors do not need a body
		if err.Error() == "" {
			req.Response.WriteHeader(err.Code())
		} else {
			if app.GetContext().Params.Debug {
				logger.Errorf("message: \"%s\", %s", err.Error(), req)
			}

			http.Error(req.Response, fmt.Sprintf(`{"error":"%s"}`, err.Error()), err.Code())
		}
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
	req.Response.WriteHeader(code)

	if obj != nil {
		if err := json.NewEncoder(req.Response).Encode(obj); err != nil {
			req.Response.WriteHeader(http.StatusInternalServerError)
			logger.Errorf("Could not write JSON response: %s", err.Error())
		}
	}
}
