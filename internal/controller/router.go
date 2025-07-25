package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/user"
	"github.com/kulikovroman08/reviewlink-backend/internal/service"
	"github.com/kulikovroman08/reviewlink-backend/pkg/middleware"
)

func SetupRouter(userService service.UserService) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	authHandler := user.NewHandler(userService)

	// Auth endpoint
	r.POST("/signup", authHandler.Signup)
	r.POST("/login", authHandler.Login)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/profile", authHandler.GetProfile)
		authorized.POST("/reviews", authHandler.CreateReview)
	}

	return r
}
