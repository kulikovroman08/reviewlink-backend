package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/pkg/middleware"
)

func SetupRouter(app *Application) *gin.Engine {
	r := gin.Default()

	// Публичные маршруты
	public := r.Group("/")
	{
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
		public.POST("/signup", app.Signup)
		public.POST("/login", app.Login)
	}

	// Защищенные маршруты
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users", app.GetUser)
		protected.PUT("/users", app.UpdateUser)
		protected.DELETE("/users", app.DeleteUser)
		protected.POST("/places", app.CreatePlace)
		protected.POST("/reviews", app.SubmitReview)
		protected.POST("/admin/tokens", app.GenerateTokens)
	}

	return r
}
