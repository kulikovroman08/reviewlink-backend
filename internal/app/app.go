package app

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller"
	repoAdmin "github.com/kulikovroman08/reviewlink-backend/internal/repository/admin"
	repoLeaderboard "github.com/kulikovroman08/reviewlink-backend/internal/repository/leaderboard"
	repoPlace "github.com/kulikovroman08/reviewlink-backend/internal/repository/place"
	repoReview "github.com/kulikovroman08/reviewlink-backend/internal/repository/review"
	repoToken "github.com/kulikovroman08/reviewlink-backend/internal/repository/token"
	repoUser "github.com/kulikovroman08/reviewlink-backend/internal/repository/user"
	svcAdmin "github.com/kulikovroman08/reviewlink-backend/internal/service/admin"
	svcLeaderboard "github.com/kulikovroman08/reviewlink-backend/internal/service/leaderboard"
	svcPlace "github.com/kulikovroman08/reviewlink-backend/internal/service/place"
	svcReview "github.com/kulikovroman08/reviewlink-backend/internal/service/review"
	svcToken "github.com/kulikovroman08/reviewlink-backend/internal/service/token"
	svcUser "github.com/kulikovroman08/reviewlink-backend/internal/service/user"
)

func InitApp(cfg *configs.Config) *gin.Engine {
	dbpool, err := pgxpool.New(context.Background(), cfg.DBUrl)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	userRepo := repoUser.NewPostgresUserRepository(dbpool)
	reviewRepo := repoReview.NewPostgresReviewRepository(dbpool)
	placeRepo := repoPlace.NewPostgresPlaceRepository(dbpool)
	tokenRepo := repoToken.NewPostgresTokenRepository(dbpool)
	adminRepo := repoAdmin.NewPostgresAdminRepository(dbpool)
	leaderboardRepo := repoLeaderboard.NewRepository(dbpool)

	tokenService := svcToken.NewTokenService(tokenRepo, cfg)
	userService := svcUser.NewUserService(userRepo)
	placeService := svcPlace.NewPlaceService(placeRepo, tokenService, cfg)
	reviewService := svcReview.NewReviewService(reviewRepo, userRepo, placeRepo, tokenService)
	adminService := svcAdmin.NewAdminService(adminRepo)
	leaderboardService := svcLeaderboard.NewService(leaderboardRepo)

	app := controller.NewApplication(userService,
		placeService,
		reviewService,
		tokenService,
		adminService,
		leaderboardService,
	)

	return controller.SetupRouter(app)
}
