package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/internal/auth"
)

func SetupRouter(userRepo auth.UserRepository) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	authHandler := auth.NewHandler(userRepo)

	// Auth endpoint
	r.POST("/signup", authHandler.Signup)
	r.POST("/login", authHandler.Login)

	authorized := r.Group("/")
	authorized.Use(auth.AuthMiddleware())
	{
		authorized.GET("/profiel", authHandler.GetProfiel)
		authorized.POST("/reviews", authHandler.CreateReview)
	}

	return r
}
