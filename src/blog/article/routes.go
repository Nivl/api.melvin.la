package article

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Nivl/api.melvin.la/src/app"
	"github.com/Nivl/api.melvin.la/src/router"
	"github.com/gorilla/mux"
)

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	routes := &router.Endpoints{
		{
			Verb:    "GET",
			Path:    "/",
			Handler: getArticles,
			Auth:    nil,
		},
		{
			Verb:    "GET",
			Path:    "/{id}",
			Handler: getArticle,
			Auth:    nil,
		},
		{
			Verb:             "POST",
			Path:             "/",
			Handler:          addArticle,
			Auth:             nil,
			JSONBodyTemplate: Article{},
		},
		{
			Verb:    "PATCH",
			Path:    "/{id}",
			Handler: updateArticle,
			Auth:    nil,
		},
	}

	routes.Activate(r)
}

func getArticles(req *router.Request) {
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

func getArticle(req *router.Request) {
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

// TODO trim data before inserting
func addArticle(req *router.Request) {
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

// TODO trim data before updating
func updateArticle(req *router.Request) {
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
