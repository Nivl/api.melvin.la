package article

import (
	"github.com/Nivl/api.melvin.la/src/app"
	"github.com/Nivl/api.melvin.la/src/router"
)

func ArticleList(req *router.Request) {
	appCtx := app.GetContext()
	doc := appCtx.DB.C("article")
	articles := []Article{}

	if err := doc.Find(nil).Sort("-createdAt").All(&articles); err != nil {
		req.ServerError(err)
		return
	}

	// httpResponse.Ok(gin, httpResponse.Collection{ToCollection(articles)})
	req.NoContent()
}
