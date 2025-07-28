package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/user"
	"github.com/kulikovroman08/reviewlink-backend/pkg/middleware"
)

func SetupRouter(app *Application) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	userHandler := user.NewHandler(app.UserService)

	r.POST("/signup", userHandler.Signup)
	r.POST("/login", userHandler.Login)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/profile", userHandler.GetProfile)
		// authorized.POST("/reviews", reviewHandler.CreateReview) — подключишь позже
	}

	return r
}
