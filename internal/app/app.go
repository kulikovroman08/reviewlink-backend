package app

import (
	"database/sql"
	"github.com/kulikovroman08/reviewlink-backend/internal/auth/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/router"
	_ "github.com/lib/pq"
)

func InitApp(cfg *configs.Config) *gin.Engine {
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	userRepo := repository.NewPostgresUserRepository(db)

	return router.SetupRouter(userRepo)
}
