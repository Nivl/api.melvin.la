package httpres

import (
	"fmt"
	"net/http"
)

// ErrorJSON replies to the request with the specified error message and HTTP code.
// The error message should be valid json.
func ErrorJSON(w http.ResponseWriter, err string, code int) {
	SetJSON(w, code)
	fmt.Fprintln(w, err)
}

// SetJSON set the response to JSON and with the specify HTTP code.
// The error message should be valid json.
func SetJSON(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
}
