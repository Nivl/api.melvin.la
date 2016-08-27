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

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func start() {
	api := gin.Default()
	api.Use(corsMiddleware())
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
