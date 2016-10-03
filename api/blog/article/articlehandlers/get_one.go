package articlehandlers

import (
	"github.com/Nivl/api.melvin.la/api/app"
	"github.com/Nivl/api.melvin.la/api/blog/article"
	"github.com/Nivl/api.melvin.la/api/router"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetOne(req *router.Request) {
	appCtx := app.GetContext()
	doc := appCtx.DB.C("article")
	article := article.Article{}
	id := bson.ObjectIdHex(mux.Vars(req.Request)["id"])

	if err := doc.Find(bson.M{"_id": id}).One(&article); err != nil {
		if err == mgo.ErrNotFound {
			req.NotFound("Article %s not found", id)
			return
		}
		req.ServerError(err)
		return
	}

	//httpResponse.Ok(gin, httpResponse.Resource{article})
	req.NoContent()
}
