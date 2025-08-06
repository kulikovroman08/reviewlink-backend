package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/pkg/middleware"
)

func SetupRouter(app *Application) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/signup", app.Signup)
	r.POST("/login", app.Login)

	authorized := r.Group("/users")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("", app.GetUser)
		authorized.PUT("", app.UpdateUser)
		authorized.DELETE("", app.DeleteUser)
	}

	return r
}
