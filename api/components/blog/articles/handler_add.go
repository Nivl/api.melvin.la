package articles

import (
	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/Nivl/api.melvin.la/api/router"
)

type HandlerAddParams struct {
	Title       string `from:"form" json:"title,omitempty" params:"required,trim"`
	Subtitle    string `from:"form" json:"subtitle,omitempty"`
	Description string `from:"form" json:"description,omitempty"`
	Content     string `from:"form" json:"content,omitempty"`
}

// HandlerAdd represents a API handler to add a new article
func HandlerAdd(req *router.Request) {
	params, ok := req.Params.(*HandlerAddParams)
	if !ok {
		req.Error(apierror.NewServerError("Couldn't cast params"))
	}

	a := &Article{
		Title:       params.Title,
		Subtitle:    params.Subtitle,
		Content:     params.Content,
		Description: params.Description,
		IsDeleted:   false,
		IsPublished: false,
	}

	// if err := a.Save(); err != nil {
	// 	req.Error(err)
	// }

	req.Created(a)
}
