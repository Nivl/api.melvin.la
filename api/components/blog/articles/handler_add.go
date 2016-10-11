package articles

import (
	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/Nivl/api.melvin.la/api/router"
)

var HandlerAddParams = router.Params{
	Query: struct{}{},
	Form: struct {
		Title       string `params:",required,trim"`
		Subtitle    string
		Description string
		Content     string
		Slug        string
	}{},
}

// HandlerAdd represents a API handler to add a new article
func HandlerAdd(req *router.Request) {
	a := &Article{
		Title:       req.Params.Get("title"),
		Subtitle:    req.Params.Get("subtitle"),
		Content:     req.Params.Get("content"),
		Description: req.Params.Get("description"),
		IsDeleted:   false,
		IsPublished: false,
	}

	if a.Title == "" {
		req.Error(apierror.NewBadRequest("Article Missing"))
	}

	if err := a.Save(); err != nil {
		req.Error(err)
	}
	req.Created(a)
}
