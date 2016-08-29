package main

import (
	"github.com/Nivl/api.melvin.la/src/app"
	"github.com/Nivl/api.melvin.la/src/blog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"net/http"
)

//func notFound(w http.ResponseWriter, r *http.Request) {
//
//}

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
	r := mux.NewRouter()
	r.Host("api.melvin.la")
	r.Host("api.melvin.loc")
	r.Headers("Content-Type", "application/json")
	blog.SetRoutes(r.PathPrefix("/blog").Subrouter())
	//router.NotFoundHandler = http.HandlerFunc(noRoutes)

	port := ":" + app.GetContext().Params.Port
	http.ListenAndServe(port, r)

	//api := gin.Default()
	//api.Use(corsMiddleware())
	//api.NoRoute(noRoute)
	//blog.SetRoutes(api)
	//api.Run(":" + app.GetContext().Params.Port)
}

func main() {
	app.InitContex()
	defer app.GetContext().Destroy()

	ensureIndexes()
	start()
}
