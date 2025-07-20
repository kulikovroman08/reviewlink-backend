package app

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/auth"
	"github.com/kulikovroman08/reviewlink-backend/internal/router"
	_ "github.com/lib/pq"
	"log"
)

func InitApp(cfg *configs.Config) *gin.Engine {
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	userRepo := auth.NewPostgresUserRepository(db)

	return router.SetupRouter(userRepo)
}
