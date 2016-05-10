package blog

import (
	"github.com/Nivl/api.melvin.la/blog/article"
	"github.com/gin-gonic/gin"
)

// SetRoutes is used to set all the routes of the blog
func SetRoutes(r *gin.Engine) {
	blogRoutes := r.Group("blog")
	article.SetRoutes(blogRoutes)
}
