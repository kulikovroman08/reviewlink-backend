package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/pkg/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

		public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		public.GET("/places/:id/reviews", app.GetReviews)

		public.GET("/leaderboard/users", app.GetUserLeaderboard)
		public.GET("/leaderboard/places", app.GetPlaceLeaderboard)
		public.GET("/leaderboard/bonuses", app.GetBonusLeaderboard)

	}

	// Защищенные маршруты
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users", app.GetUser)
		protected.PUT("/users", app.UpdateUser)
		protected.DELETE("/users", app.DeleteUser)
		protected.GET("/users/stats", app.GetUserStats)

		protected.POST("/places", app.CreatePlace)

		protected.POST("/reviews", app.SubmitReview)
		protected.PATCH("/reviews/:id", app.UpdateReview)
		protected.DELETE("/reviews/:id", app.DeleteReview)

		protected.POST("/admin/tokens", app.GenerateTokens)
		protected.GET("/admin/stats", app.GetStats)

		protected.POST("/bonuses/redeem", app.RedeemBonus)
		protected.GET("/bonuses", app.GetUserBonuses)
		protected.POST("/bonuses/validate", app.ValidateBonus)
	}

	return r
}
