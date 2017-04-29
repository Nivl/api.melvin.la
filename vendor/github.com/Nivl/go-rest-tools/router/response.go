package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Nivl/go-rest-tools/logger"
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/network/http/httpres"
	"github.com/Nivl/go-rest-tools/notifiers/mailer"
)

// Error sends an error to the client
// If the error is an instance of HTTPError, the returned code will
// match HTTPError.Code(). It returns a 500 if no code has been set.
func (req *Request) Error(e error) {
	if req == nil {
		return
	}

	err, casted := e.(*httperr.HTTPError)
	if !casted {
		err = httperr.NewServerError(e.Error()).(*httperr.HTTPError)
	}

	switch err.Code() {
	case http.StatusInternalServerError:
		httpres.ErrorJSON(req.Response, `{"error":"Something went wrong"}`, http.StatusInternalServerError)
	default:
		// Some errors do not need a body
		if err.Error() == "" {
			req.Response.WriteHeader(err.Code())
		} else {
			httpres.ErrorJSON(req.Response, fmt.Sprintf(`{"error":"%s"}`, err.Error()), err.Code())
		}
	}

	logger.Errorf(`code: "%d", message: "%s", %s`, err.Code(), err.Error(), req)

	// We send an email for all server error
	if err.Code() == http.StatusInternalServerError {

		if mailer.Emailer != nil {
			sendEmail := func(stacktrace []byte) {
				err := mailer.Emailer.SendStackTrace(stacktrace, req.Endpoint(), err.Error(), req.ID)
				if err != nil {
					logger.Error(err.Error())
				}
			}

			go sendEmail(debug.Stack())
		}
	}
}

// NoContent sends a http.StatusNoContent response
// It should be used for successful DELETE requests
func (req *Request) NoContent() {
	if req == nil {
		return
	}

	req.Response.WriteHeader(http.StatusNoContent)
}

// Created sends a http.StatusCreated response with a JSON object attached
// It should be used for successful POST requests
func (req *Request) Created(obj interface{}) {
	if req == nil {
		return
	}

	req.RenderJSON(http.StatusCreated, obj)
}

// Ok sends a http.StatusOK response with a JSON object attached
// It should be used for successful GET, PATCH, and PUR requests
func (req *Request) Ok(obj interface{}) {
	if req == nil {
		return
	}

	req.RenderJSON(http.StatusOK, obj)
}

// RenderJSON attaches a json object to the response
func (req *Request) RenderJSON(code int, obj interface{}) {
	httpres.SetJSON(req.Response, code)

	if obj != nil {
		if err := json.NewEncoder(req.Response).Encode(obj); err != nil {
			req.Error(fmt.Errorf("Could not write JSON response: %s", err.Error()))
		}
	}
}
