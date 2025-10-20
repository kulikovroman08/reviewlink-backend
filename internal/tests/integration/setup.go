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

	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository/place"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/kulikovroman08/reviewlink-backend/internal/controller"
	reviewRepo "github.com/kulikovroman08/reviewlink-backend/internal/repository/review"
	tokenRepo "github.com/kulikovroman08/reviewlink-backend/internal/repository/token"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository/user"
	placeService "github.com/kulikovroman08/reviewlink-backend/internal/service/place"
	reviewService "github.com/kulikovroman08/reviewlink-backend/internal/service/review"
	tokenService "github.com/kulikovroman08/reviewlink-backend/internal/service/token"
	userService "github.com/kulikovroman08/reviewlink-backend/internal/service/user"
)

type TestSetup struct {
	App *gin.Engine
	DB  *pgxpool.Pool
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

	tokSrv := tokenService.NewTokenService(tokRepo, &cfg)
	userSrv := userService.NewUserService(userRepo)
	placeSrv := placeService.NewPlaceService(placeRepo, tokSrv, &cfg)
	reviewSrv := reviewService.NewReviewService(reviewRepo, userRepo, placeRepo, tokSrv)

	app := controller.NewApplication(userSrv, placeSrv, reviewSrv, tokSrv)
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
