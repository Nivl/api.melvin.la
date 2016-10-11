package articles

import "github.com/Nivl/api.melvin.la/api/router"

// HandlerGet represents a API handler to get a single article
func HandlerGet(req *router.Request) {
	// appCtx := app.GetContext()
	// doc := appCtx.DB.C("article")
	// article := article.Article{}
	// id := bson.ObjectIdHex(mux.Vars(req.Request)["id"])

	// if err := doc.Find(bson.M{"_id": id}).One(&article); err != nil {
	// 	if err == mgo.ErrNotFound {
	// 		req.NotFound("Article %s not found", id)
	// 		return
	// 	}
	// 	req.Error(NewSer)
	// 	return
	// }

	//httpResponse.Ok(gin, httpResponse.Resource{article})
	req.NoContent()
}
