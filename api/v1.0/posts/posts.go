package posts

import (
	"github.com/gin-gonic/gin"
	"loundry_rest/lib/middlewares"
)

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	posts := r.Group("/posts")
	{
		posts.POST("/", middlewares.Authorized, create)
		posts.GET("/", middlewares.Authorized, list)
		posts.GET("/getPostById", middlewares.Authorized, readByParams)
		posts.GET("/readById/:id", middlewares.Authorized, readById)
		posts.DELETE("/:id", middlewares.Authorized, remove)
		posts.POST("/ubahData", middlewares.Authorized, update)
	}
}
