package app

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller"
	repoPlace "github.com/kulikovroman08/reviewlink-backend/internal/repository/place"
	repoReview "github.com/kulikovroman08/reviewlink-backend/internal/repository/review"
	repoUser "github.com/kulikovroman08/reviewlink-backend/internal/repository/user"
	servicePlace "github.com/kulikovroman08/reviewlink-backend/internal/service/place"
	serviceReview "github.com/kulikovroman08/reviewlink-backend/internal/service/review"
	serviceUser "github.com/kulikovroman08/reviewlink-backend/internal/service/user"
)

func InitApp(cfg *configs.Config) *gin.Engine {
	dbpool, err := pgxpool.New(context.Background(), cfg.DBUrl)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	userRepo := repoUser.NewPostgresUserRepository(dbpool)
	userService := serviceUser.NewService(userRepo)

	placeRepo := repoPlace.NewPostgresPlaceRepository(dbpool)
	placeService := servicePlace.NewService(placeRepo)

	reviewRepo := repoReview.NewPostgresReviewRepository(dbpool)
	reviewService := serviceReview.NewService(reviewRepo, userRepo)

	app := controller.NewApplication(userService, placeService, reviewService)

	return controller.SetupRouter(app)
}
