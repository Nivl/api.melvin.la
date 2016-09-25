package router

import (
	"mime/multipart"
	"net/url"
	"os"

	"bitbucket.org/splice/api/splice/logger"
)

type Params struct {
	url.Values // Contains all kind of params (qs and forms)

	Query url.Values
	Form  url.Values

	Files    map[string][]*multipart.FileHeader // Files uploaded in a multipart form
	tmpFiles []*os.File                         // Temp files used during the request.
}

func (r *Request) ParseParams() error {
	r.Params.Query = r.Request.URL.Query()

	// Parse the body depending on the content type.
	switch r.GetContentType() {
	case ctFormData:
		if err := r.Request.ParseForm(); err != nil {
			logger.Errorf("Error parsing request body %s", err.Error())
			return err
		} else {
			r.Params.Form = r.Request.Form
		}
	case ctMultipartFormData:
		twoMB := int64(2 << 20)
		if err := r.Request.ParseMultipartForm(twoMB); err != nil {
			logger.Errorf("Error parsing multipart request body - %s", err)
			return err
		} else {
			r.Params.Form = r.Request.MultipartForm.Value
			r.Params.Files = r.Request.MultipartForm.File
		}
	}

	// merge the params
	nbParams := len(r.Params.Query) + len(r.Params.Form)
	r.Params.Values = make(url.Values, nbParams)

	for k, v := range r.Params.Query {
		r.Params.Values[k] = append(r.Params.Values[k], v...)
	}
	for k, v := range r.Params.Form {
		r.Params.Values[k] = append(r.Params.Values[k], v...)
	}
	return nil
}

func removeParams(r *Request) error {
	if r.Request.MultipartForm != nil {
		err := r.Request.MultipartForm.RemoveAll()
		if err != nil {
			logger.Errorf("Error removing temporary files:" + err.Error())
			return err
		}
	}

	for _, tmpFile := range r.Params.tmpFiles {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			logger.Errorf("Could not remove upload temp file:" + err.Error())
			return err
		}
	}

	return nil
}
