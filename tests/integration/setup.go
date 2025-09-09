package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"time"

	"github.com/kulikovroman08/reviewlink-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/kulikovroman08/reviewlink-backend/internal/controller"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository/user"
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

	app := &controller.Application{
		UserService: userSrv,
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
	}

	return &TestSetup{
		App: r,
		DB:  db,
	}
}

func (ts *TestSetup) TruncateUsers() {
	_, err := ts.DB.Exec(context.Background(), "TRUNCATE users RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatalf("failed to truncate users table: %v", err)
	}
}

func (ts *TestSetup) SignupAndLogin(email, password string) string {
	payload := map[string]string{
		"name":     "Test User",
		"email":    email,
		"password": password,
	}
	data, _ := json.Marshal(payload)

	// Signup
	req := httptest.NewRequest("POST", "/signup", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ts.App.ServeHTTP(rec, req)

	// Login
	req = httptest.NewRequest("POST", "/login", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	ts.App.ServeHTTP(rec, req)

	var resp map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&resp)

	return resp["token"]
}
