package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/app"
	"github.com/melvin-laplanche/ml-api/src/logger"
	"github.com/melvin-laplanche/ml-api/src/mlhttp"
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
		mlhttp.ErrorJSON(req.Response, `{"error":"Something went wrong"}`, http.StatusInternalServerError)
	default:
		// Some errors do not need a body
		if err.Error() == "" {
			req.Response.WriteHeader(err.Code())
		} else {
			mlhttp.ErrorJSON(req.Response, fmt.Sprintf(`{"error":"%s"}`, err.Error()), err.Code())
		}
	}

	logger.Errorf(`code: "%d", message: "%s", %s`, err.Code(), err.Error(), req)

	// We send an email for all server error
	if err.Code() == 500 {
		context := app.GetContext()

		if context != nil && context.Mailer != nil {
			sendEmail := func(stacktrace []byte) {
				err := context.Mailer.SendStackTrace(stacktrace, req.Endpoint(), err.Error(), req.ID)
				if err != nil {
					logger.Error(err.Error())
				}
			}

			go sendEmail(debug.Stack())
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
	mlhttp.SetJSON(req.Response, code)

	if obj != nil {
		if err := json.NewEncoder(req.Response).Encode(obj); err != nil {
			req.Error(fmt.Errorf("Could not write JSON response: %s", err.Error()))
		}
	}
}
