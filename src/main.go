package main

import (
	"github.com/Nivl/api.melvin.la/src/app"
	"github.com/Nivl/api.melvin.la/src/blog"
	"github.com/Nivl/api.melvin.la/src/http-response"
	"github.com/gin-gonic/gin"
)

func noRoute(gin *gin.Context) {
	httpResponse.NotFound(gin)
}

func ensureIndexes() {
	blog.EnsureIndexes()
}

func start() {
	api := gin.Default()
	api.NoRoute(noRoute)
	blog.SetRoutes(api)
	api.Run(":" + app.GetContext().Params.Port)
}

func main() {
	app.InitContex()
	defer app.GetContext().Destroy()

	ensureIndexes()
	start()
}
