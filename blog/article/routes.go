package article

import (
	"log"

	"github.com/Nivl/api.melvin.la/app"
	"github.com/Nivl/api.melvin.la/http-response"
	"github.com/gin-gonic/gin"
)

// SetRoutes is used to set all the routes of the article
func SetRoutes(blog *gin.RouterGroup) {
	articles := blog.Group("articles")
	articles.GET("/", getArticles)
}

func getArticles(gin *gin.Context) {
	appCtx := app.GetContext()
	doc := appCtx.DB.C("article")
	articles := []Article{}

	err := doc.Find(nil).Sort("-createdAt").All(&articles)

	if err != nil {
		log.Println(err.Error())
		httpResponse.ServerError(gin)
	} else {
		httpResponse.Ok(gin, httpResponse.Collection{ToCollection(articles)})
	}
}
