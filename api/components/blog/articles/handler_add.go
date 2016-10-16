package articles

import "github.com/Nivl/api.melvin.la/api/router"

type HandlerAddParams struct {
	Title       string `from:"query" json:"title" params:"required,trim"`
	Subtitle    string `from:"query" json:"subtitle,omitempty"`
	Description string `from:"query" json:"description,omitempty"`
	Content     string `from:"query" json:"content,omitempty"`
}

// HandlerAdd represents a API handler to add a new article
func HandlerAdd(req *router.Request) {

	params := req.Params.(HandlerAddParams)

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
