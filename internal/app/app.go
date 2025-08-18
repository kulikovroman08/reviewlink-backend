package app

import (
	"context"
	"log"

	"github.com/kulikovroman08/reviewlink-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller"
	repoPlace "github.com/kulikovroman08/reviewlink-backend/internal/repository/place"
	repoReview "github.com/kulikovroman08/reviewlink-backend/internal/repository/review"
	repoUser "github.com/kulikovroman08/reviewlink-backend/internal/repository/user"
)

func InitApp(cfg *configs.Config) *gin.Engine {
	dbpool, err := pgxpool.New(context.Background(), cfg.DBUrl)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	userRepo := repoUser.NewPostgresUserRepository(dbpool)
	reviewRepo := repoReview.NewPostgresReviewRepository(dbpool)
	placeRepo := repoPlace.NewPostgresPlaceRepository(dbpool)

	svc := service.NewService(userRepo, reviewRepo, placeRepo)

	app := controller.NewApplication(svc)

	return controller.SetupRouter(app)
}
