package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"

	"github.com/kulikovroman08/reviewlink-backend/internal/controller"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository/user"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/configs"
	_ "github.com/lib/pq"
)

func InitApp(cfg *configs.Config) *gin.Engine {
	dbpool, err := pgxpool.New(context.Background(), cfg.DBUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	userRepo := user.NewPostgresUserRepository(dbpool)

	return controller.SetupRouter(userRepo)
}
