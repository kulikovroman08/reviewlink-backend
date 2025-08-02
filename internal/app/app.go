package app

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller"
	repoUser "github.com/kulikovroman08/reviewlink-backend/internal/repository/user"
	serviceUser "github.com/kulikovroman08/reviewlink-backend/internal/service/user"
)

func InitApp(cfg *configs.Config) *gin.Engine {
	dbpool, err := pgxpool.New(context.Background(), cfg.DBUrl)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	userRepo := repoUser.NewPostgresUserRepository(dbpool)
	userService := serviceUser.NewService(userRepo)

	app := controller.NewApplication(userService)

	return controller.SetupRouter(app)
}
