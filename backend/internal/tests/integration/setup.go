package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"

	repoVote "github.com/kulikovroman08/reviewlink-backend/internal/repository/review_vote"
	svcVote "github.com/kulikovroman08/reviewlink-backend/internal/service/review_vote"

	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository/place"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/kulikovroman08/reviewlink-backend/internal/controller"
	osm "github.com/kulikovroman08/reviewlink-backend/internal/infra/client/osm"
	repoAdmin "github.com/kulikovroman08/reviewlink-backend/internal/repository/admin"
	bonusRepo "github.com/kulikovroman08/reviewlink-backend/internal/repository/bonus"
	repoLeaderboard "github.com/kulikovroman08/reviewlink-backend/internal/repository/leaderboard"
	repoOwnerRequest "github.com/kulikovroman08/reviewlink-backend/internal/repository/owner_request"
	restrictionRepo "github.com/kulikovroman08/reviewlink-backend/internal/repository/restriction"
	reviewRepo "github.com/kulikovroman08/reviewlink-backend/internal/repository/review"
	repoReply "github.com/kulikovroman08/reviewlink-backend/internal/repository/review_reply"
	tokenRepo "github.com/kulikovroman08/reviewlink-backend/internal/repository/token"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository/user"
	adminService "github.com/kulikovroman08/reviewlink-backend/internal/service/admin"
	svcBonus "github.com/kulikovroman08/reviewlink-backend/internal/service/bonus"
	svcLeaderboard "github.com/kulikovroman08/reviewlink-backend/internal/service/leaderboard"
	svcOwnerRequest "github.com/kulikovroman08/reviewlink-backend/internal/service/owner_request"
	svcPlace "github.com/kulikovroman08/reviewlink-backend/internal/service/place"
	reviewService "github.com/kulikovroman08/reviewlink-backend/internal/service/review"
	svcReply "github.com/kulikovroman08/reviewlink-backend/internal/service/review_reply"
	tokenService "github.com/kulikovroman08/reviewlink-backend/internal/service/token"
	userService "github.com/kulikovroman08/reviewlink-backend/internal/service/user"
)

type TestSetup struct {
	App *gin.Engine
	DB  *pgxpool.Pool
}

type stubOSMClient struct{}

func (s stubOSMClient) SearchPlaces(ctx context.Context, p osm.SearchParams) ([]osm.Place, error) {
	return []osm.Place{}, nil
}

type stubPublicNewsService struct{}

func (s stubPublicNewsService) ListNews(ctx context.Context, limit int) ([]model.NewsItem, *time.Time, error) {
	return []model.NewsItem{}, nil, nil
}

func NewTestSetup() *TestSetup {
	gin.SetMode(gin.TestMode)

	root := os.Getenv("PROJECT_ROOT")
	if root == "" {
		root, _ = os.Getwd()
	}
	_ = godotenv.Load(filepath.Join(root, ".env.test"))

	cfg := configs.LoadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbURL := os.Getenv("DB_URL_TEST")
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to test DB: %v", err)
	}

	userRepo := user.NewPostgresUserRepository(db)
	placeRepo := place.NewPostgresPlaceRepository(db)
	reviewRepo := reviewRepo.NewPostgresReviewRepository(db)
	tokRepo := tokenRepo.NewPostgresTokenRepository(db)
	adminRepo := repoAdmin.NewPostgresAdminRepository(db)
	leaderboardRepo := repoLeaderboard.NewRepository(db)
	bonusRepo := bonusRepo.NewPostgresBonusRepository(db)
	restrictionRepo := restrictionRepo.NewPostgresUserRestrictionRepository(db)
	replyRepo := repoReply.NewPostgresReviewReplyRepository(db)
	voteRepo := repoVote.NewPostgresReviewVoteRepository()
	ownerRequestRepo := repoOwnerRequest.NewPostgresOwnerRequestRepository(db)

	tokSrv := tokenService.NewTokenService(tokRepo, placeRepo, &cfg)
	userSrv := userService.NewUserService(userRepo, reviewRepo, bonusRepo)
	osmClient := stubOSMClient{}
	placeSrv := svcPlace.NewPlaceService(placeRepo, tokSrv, osmClient, &cfg)
	reviewSrv := reviewService.NewReviewService(reviewRepo, userRepo, placeRepo, restrictionRepo, replyRepo)
	adminSrv := adminService.NewAdminService(adminRepo)
	leaderboardSrv := svcLeaderboard.NewService(leaderboardRepo, nil)
	bonusSrv := svcBonus.NewBonusService(userRepo, bonusRepo, &cfg)
	reviewReplySrv := svcReply.NewReviewReplyService(reviewRepo, placeRepo, replyRepo)
	reviewVoteService := svcVote.NewReviewVoteService(db, voteRepo, reviewRepo, userRepo)
	publicNewsService := stubPublicNewsService{}
	ownerRequestService := svcOwnerRequest.NewOwnerRequestService(ownerRequestRepo, placeRepo, userRepo)

	app := controller.NewApplication(
		userSrv,
		placeSrv,
		ownerRequestService,
		publicNewsService,
		reviewSrv,
		tokSrv,
		adminSrv,
		leaderboardSrv,
		bonusSrv,
		reviewReplySrv,
		reviewVoteService,
	)

	r := controller.SetupRouter(app)

	return &TestSetup{
		App: r,
		DB:  db,
	}
}

func (ts *TestSetup) Close() {
	ts.DB.Close()
}

func (ts *TestSetup) Login(email, password string) string {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	data, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ts.App.ServeHTTP(rec, req)

	var resp map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&resp)

	return resp["token"]
}
