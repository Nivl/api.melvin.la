package testhelpers

import (
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/Nivl/api.melvin.la/api/components/api"
	"github.com/Nivl/api.melvin.la/api/router"
)

func NewRequest(e router.Endpoint, uri string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(e.Verb, uri, nil)
	if err != nil {
		log.Fatal(err)
	}

	rec := httptest.NewRecorder()
	r := api.GetRouter()
	r.ServeHTTP(rec, req)
	return rec
}
