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

	// Places
	places := r.Group("/places")
	places.Use(middleware.AuthMiddleware())
	{
		places.POST("", app.CreatePlace)
	}

	// Reviews
	reviews := r.Group("/reviews")
	reviews.Use(middleware.AuthMiddleware())
	{
		reviews.POST("", app.SubmitReview)
	}

	// Generate Token Admin
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	{
		admin.POST("/tokens", app.GenerateTokens)
	}

	return r
}
