package main

import (
	"github.com/Nivl/api.melvin.la/app"
	"github.com/Nivl/api.melvin.la/blog"
	"github.com/gin-gonic/gin"
)

func main() {
	appContext := app.GetContext()
	defer appContext.Destroy()

	api := gin.Default()
	blog.SetRoutes(api)
	api.Run(":" + appContext.Params.Port)
}
