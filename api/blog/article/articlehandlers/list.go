package articlehandlers

import (
	"github.com/Nivl/api.melvin.la/api/blog/article"
	"github.com/Nivl/api.melvin.la/api/router"
)

// List returns a list of Articles
func List(req *router.Request) {
	articles := []*article.Article{}

	if err := article.Query().Find(nil).Sort("-createdAt").All(&articles); err != nil {
		req.ServerError(err)
		return
	}

	req.Ok(article.NewPayloadFromModels(articles))
}
