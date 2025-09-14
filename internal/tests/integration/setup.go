package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"time"

	"github.com/kulikovroman08/reviewlink-backend/internal/repository/place"

	"github.com/kulikovroman08/reviewlink-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/kulikovroman08/reviewlink-backend/internal/controller"
	reviewRepo "github.com/kulikovroman08/reviewlink-backend/internal/repository/review"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository/user"
	placeService "github.com/kulikovroman08/reviewlink-backend/internal/service/place"
	reviewService "github.com/kulikovroman08/reviewlink-backend/internal/service/review"
	userService "github.com/kulikovroman08/reviewlink-backend/internal/service/user"
)

type TestSetup struct {
	App *gin.Engine
	DB  *pgxpool.Pool
}

func NewTestSetup() *TestSetup {
	gin.SetMode(gin.TestMode)

	_ = godotenv.Load(".env.test")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbURL := os.Getenv("DB_URL")
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to test DB: %v", err)
	}

	userRepo := user.NewPostgresUserRepository(db)
	userSrv := userService.NewUserService(userRepo)
	placeRepo := place.NewPostgresPlaceRepository(db)
	placeSrv := placeService.NewPlaceService(placeRepo)
	reviewRepo := reviewRepo.NewPostgresReviewRepository(db)
	reviewSrv := reviewService.NewReviewService(reviewRepo, userRepo, placeRepo)
	//tokRepo := tokenRepo.NewPostgresTokenRepository(db)
	//tokSrv := tokenService.NewTokenService(tokRepo)

	app := &controller.Application{
		UserService:   userSrv,
		PlaceService:  placeSrv,
		ReviewService: reviewSrv,
		//TokenService:  tokSrv,
	}

	r := gin.Default()
	r.POST("/signup", app.Signup)
	r.POST("/login", app.Login)

	// Защищённые маршруты
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users", app.GetUser)
		protected.PUT("/users", app.UpdateUser)
		protected.DELETE("/users", app.DeleteUser)
		protected.POST("/places", app.CreatePlace)
		protected.POST("/reviews", app.SubmitReview)
		//protected.POST("/admin/tokens", app.GenerateTokens)
	}

	return &TestSetup{
		App: r,
		DB:  db,
	}
}

func (ts *TestSetup) TruncateAll() {
	_, err := ts.DB.Exec(context.Background(),
		"TRUNCATE users, places, reviews, review_tokens RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatalf("failed to truncate users table: %v", err)
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
