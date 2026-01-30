package app

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller"
	repoAdmin "github.com/kulikovroman08/reviewlink-backend/internal/repository/admin"
	bonusRepo "github.com/kulikovroman08/reviewlink-backend/internal/repository/bonus"
	repoLeaderboard "github.com/kulikovroman08/reviewlink-backend/internal/repository/leaderboard"
	repoPlace "github.com/kulikovroman08/reviewlink-backend/internal/repository/place"
	restrictionRepo "github.com/kulikovroman08/reviewlink-backend/internal/repository/restriction"
	repoReview "github.com/kulikovroman08/reviewlink-backend/internal/repository/review"
	repoReply "github.com/kulikovroman08/reviewlink-backend/internal/repository/review_reply"
	repoVote "github.com/kulikovroman08/reviewlink-backend/internal/repository/review_vote"
	repoToken "github.com/kulikovroman08/reviewlink-backend/internal/repository/token"
	repoUser "github.com/kulikovroman08/reviewlink-backend/internal/repository/user"
	svcAdmin "github.com/kulikovroman08/reviewlink-backend/internal/service/admin"
	svcBonus "github.com/kulikovroman08/reviewlink-backend/internal/service/bonus"
	svcLeaderboard "github.com/kulikovroman08/reviewlink-backend/internal/service/leaderboard"
	svcNews "github.com/kulikovroman08/reviewlink-backend/internal/service/news"
	svcPlace "github.com/kulikovroman08/reviewlink-backend/internal/service/place"
	svcReview "github.com/kulikovroman08/reviewlink-backend/internal/service/review"
	svcReply "github.com/kulikovroman08/reviewlink-backend/internal/service/review_reply"
	svcVote "github.com/kulikovroman08/reviewlink-backend/internal/service/review_vote"
	svcToken "github.com/kulikovroman08/reviewlink-backend/internal/service/token"
	svcUser "github.com/kulikovroman08/reviewlink-backend/internal/service/user"
)

func InitApp(cfg *configs.Config) *gin.Engine {
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	rdb := initRedis(ctx, cfg)

	var lbCache svcLeaderboard.LeaderboardCache
	if rdb != nil {
		lbCache = svcLeaderboard.NewRedisLeaderboardCache(rdb, 60*time.Second)
	}

	userRepo := repoUser.NewPostgresUserRepository(dbpool)
	reviewRepo := repoReview.NewPostgresReviewRepository(dbpool)
	placeRepo := repoPlace.NewPostgresPlaceRepository(dbpool)
	tokenRepo := repoToken.NewPostgresTokenRepository(dbpool)
	adminRepo := repoAdmin.NewPostgresAdminRepository(dbpool)
	leaderboardRepo := repoLeaderboard.NewRepository(dbpool)
	bonusRepo := bonusRepo.NewPostgresBonusRepository(dbpool)
	restrictionRepo := restrictionRepo.NewPostgresUserRestrictionRepository(dbpool)
	replyRepo := repoReply.NewPostgresReviewReplyRepository(dbpool)
	voteRepo := repoVote.NewPostgresReviewVoteRepository()

	tokenService := svcToken.NewTokenService(tokenRepo, placeRepo, cfg)
	userService := svcUser.NewUserService(userRepo, reviewRepo, bonusRepo)
	osmClient := buildOSMClient(cfg, rdb)
	placeService := svcPlace.NewPlaceService(placeRepo, tokenService, osmClient, cfg)
	reviewService := svcReview.NewReviewService(reviewRepo, userRepo, placeRepo, restrictionRepo, replyRepo)
	adminService := svcAdmin.NewAdminService(adminRepo)
	leaderboardService := svcLeaderboard.NewService(leaderboardRepo, lbCache)
	bonusService := svcBonus.NewBonusService(userRepo, bonusRepo, cfg)
	reviewReplyService := svcReply.NewReviewReplyService(reviewRepo, placeRepo, replyRepo)
	reviewVoteService := svcVote.NewReviewVoteService(dbpool, voteRepo, reviewRepo, userRepo)
	publicNewsClient := buildPublicNewsClient(cfg, rdb)
	publicNewsService := svcNews.NewNewsService(publicNewsClient)

	app := controller.NewApplication(
		userService,
		placeService,
		publicNewsService,
		reviewService,
		tokenService,
		adminService,
		leaderboardService,
		bonusService,
		reviewReplyService,
		reviewVoteService,
	)

	return controller.SetupRouter(app)
}

func initRedis(ctx context.Context, cfg *configs.Config) *redis.Client {
	if cfg.RedisAddr == "" {
		log.Printf("redis disabled (REDIS_ADDR is empty)")
		return nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		log.Printf("redis disabled (can't connect to %s): %v", cfg.RedisAddr, err)
		return nil
	}

	log.Printf("redis enabled: %s", cfg.RedisAddr)
	return rdb
}
