package app

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/auth"
	_ "github.com/lib/pq"
	"log"
)

func InitApp(cfg *configs.Config) *gin.Engine {
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	userRepo := auth.NewPostgresUserRepository(db)
	authHandler := auth.NewHandler(userRepo)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.POST("/signup", authHandler.Signup)

	return router
}
