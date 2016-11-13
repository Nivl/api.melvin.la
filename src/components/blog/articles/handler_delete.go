package articles

import "github.com/melvin-laplanche/ml-api/src/router"

// HandlerDelete represents a API handler to delete an article
func HandlerDelete(req *router.Request) {
	//appCtx := app.GetContext()
	//doc := appCtx.DB.C("article")
	//article := Article{
	//	ID:        bson.NewObjectId(),
	//	CreatedAt: time.Now(),
	//}
	//
	//if err := gin.Bind(&article); err != nil {
	//	req.BadRequest(err)
	//	return
	//}
	//
	//if err := doc.Insert(article); err != nil {
	//	if mgo.IsDup(err) {
	//		req.Conflict("The slug %s already exists in the database", article.Slug)
	//		return
	//	}
	//
	//	req.ServerError(err)
	//	return
	//}

	//httpResponse.Ok(gin, httpResponse.Resource{article})
	req.NoContent()
}
