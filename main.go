package main

import (
	"github.com/Nivl/api.melvin.la/app"
	"github.com/Nivl/api.melvin.la/blog"
	"github.com/Nivl/api.melvin.la/http-response"
	"github.com/gin-gonic/gin"
)

func noRoute(gin *gin.Context) {
	httpResponse.NotFound(gin)
}

func ensureIndexes() {
	blog.EnsureIndexes()
}

func main() {
	appContext := app.GetContext()
	defer appContext.Destroy()

	ensureIndexes()

	api := gin.Default()
	api.NoRoute(noRoute)
	blog.SetRoutes(api)
	api.Run(":" + appContext.Params.Port)
}
