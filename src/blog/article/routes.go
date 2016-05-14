package article

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Nivl/api.melvin.la/src/app"
	"github.com/Nivl/api.melvin.la/src/http-response"
	"github.com/gin-gonic/gin"
)

// SetRoutes is used to set all the routes of the article
func SetRoutes(blog *gin.RouterGroup) {
	articles := blog.Group("articles")
	articles.GET("/:id", getArticle)
	articles.GET("", getArticles)
	articles.POST("", addArticle)
}

func getArticles(gin *gin.Context) {
	appCtx := app.GetContext()
	doc := appCtx.DB.C("article")
	articles := []Article{}

	if err := doc.Find(nil).Sort("-createdAt").All(&articles); err != nil {
		log.Println(err.Error())
		httpResponse.ServerError(gin)
	} else {
		httpResponse.Ok(gin, httpResponse.Collection{ToCollection(articles)})
	}
}

func getArticle(gin *gin.Context) {
	appCtx := app.GetContext()
	doc := appCtx.DB.C("article")
	article := Article{}
	id := bson.ObjectIdHex(gin.Param("id"))

	if err := doc.Find(bson.M{"_id": id}).One(&article); err != nil {
		if err == mgo.ErrNotFound {
			httpResponse.NotFound(gin)
		} else {
			log.Println(err.Error())
			httpResponse.ServerError(gin)
		}
	} else {
		httpResponse.Ok(gin, httpResponse.Resource{article})
	}
}

// TODO trim data before inserting
func addArticle(gin *gin.Context) {
	appCtx := app.GetContext()
	doc := appCtx.DB.C("article")
	article := Article{
		ID:        bson.NewObjectId(),
		CreatedAt: time.Now(),
	}

	if err := gin.Bind(&article); err != nil {
		httpResponse.BadRequest(gin, err.Error())
	} else if err := doc.Insert(article); err != nil {
		if mgo.IsDup(err) {
			httpResponse.Conflict(gin, "This slug already exists in the database")
		} else {
			log.Println(err.Error())
			httpResponse.ServerError(gin)
		}
	} else {
		httpResponse.Ok(gin, httpResponse.Resource{article})
	}
}
