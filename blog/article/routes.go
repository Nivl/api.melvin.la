package article

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Nivl/api.melvin.la/app"
	"github.com/Nivl/api.melvin.la/http-response"
	"github.com/gin-gonic/gin"
)

// SetRoutes is used to set all the routes of the article
func SetRoutes(blog *gin.RouterGroup) {
	articles := blog.Group("articles")
	articles.GET("/:id", getArticle)
	articles.GET("", getArticles)
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
