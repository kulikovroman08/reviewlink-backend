package app

import (
	"database/sql"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository/auth"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/configs"
	_ "github.com/lib/pq"
)

func InitApp(cfg *configs.Config) *gin.Engine {
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	userRepo := auth.NewPostgresUserRepository(db)

	return controller.SetupRouter(userRepo)
}
