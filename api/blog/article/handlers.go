package article

import (
	"github.com/Nivl/api.melvin.la/api/app"
	"github.com/Nivl/api.melvin.la/api/router"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func List(req *router.Request) {
	articles := []*Article{}

	if err := Query().Find(nil).Sort("-createdAt").All(&articles); err != nil {
		req.ServerError(err)
		return
	}

	req.Ok(NewPayloadFromModels(articles))
}

// TODO trim data before updating
func Update(req *router.Request) {
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

func GetOne(req *router.Request) {
	appCtx := app.GetContext()
	doc := appCtx.DB.C("article")
	article := Article{}
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

func Add(req *router.Request) {
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
